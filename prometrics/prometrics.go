/*
 * Copyright 2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prometrics

import (
	"errors"
	"reflect"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"stash.kopano.io/kgol/ksurveyclient-go"
)

// DefaultRegistry exposes the registry which is used by autosurvey.
var DefaultRegistry *ksurveyclient.Registry

func init() {
	DefaultRegistry = ksurveyclient.DefaultRegistry
}

type collector struct {
	collector prometheus.Collector
	whitelist map[string]string
}

// WrapCollector takes a prometheus.Collector with a whitelist and returns a
// kcsurveyclient.Collector. If the provided whitelist is nil, all the metrics
// in the provided prometheus.Collector whill be collected. If whitelist is not
// nil, only the metrics found as keys will be collected. If the values of the
// whitelist keys is not empty, the metrics name will also be changed to that
// value.
func WrapCollector(c prometheus.Collector, whitelist map[string]string) ksurveyclient.Collector {
	return &collector{
		collector: c,
		whitelist: whitelist,
	}
}

func (c *collector) Collect(ch chan<- ksurveyclient.Metric) {
	metrics := make(chan prometheus.Metric)
	go func() {
		c.collector.Collect(metrics)
		close(metrics)
	}()
	for m := range metrics {
		pm := newProMetrics(m, c.whitelist)
		if pm != nil {
			ch <- pm
		}
	}
}

type proMetrics struct {
	metric prometheus.Metric

	fqName string
	help   string

	mtype string
	mmode string

	err error
}

func newProMetrics(metric prometheus.Metric, whitelist map[string]string) *proMetrics {
	var err error

	desc := metric.Desc()
	// NOTE(longsleep): Since all fields in Desc are private we unfortunately
	// need to use reflect to find the name and description of the prometheus
	// metrics currectly processed :(.
	d := reflect.ValueOf(*desc)
	fqNameValue := d.FieldByName("fqName")
	if !fqNameValue.IsValid() {
		err = errors.New("no fqName field")
	}

	mtype := ""
	mmode := ""
	fqName := fqNameValue.String()
	if whitelist != nil {
		alias, ok := whitelist[fqName]
		if !ok {
			return nil
		}
		if alias != "" {
			// Alias format supported `alias,mmode:mtype`.
			aliasParts := strings.SplitN(alias, ",", 2)
			alias = aliasParts[0]
			if len(aliasParts) > 1 {
				formatParts := strings.SplitN(aliasParts[1], ":", 2)
				mmode = formatParts[0]
				if len(formatParts) > 1 {
					mtype = formatParts[1]
				}
			}
			fqName = alias
		}
	}
	help := d.FieldByName("help")
	if !help.IsValid() {
		err = errors.New("no help field")
	}

	return &proMetrics{
		metric: metric,

		fqName: fqName,
		help:   help.String(),

		mtype: mtype,
		mmode: mmode,

		err: err,
	}
}

func (pm *proMetrics) Write(md *ksurveyclient.MetricData) error {
	if pm.err != nil {
		return pm.err
	}

	dtoMetric := &dto.Metric{}
	if err := pm.metric.Write(dtoMetric); err != nil {
		return err
	}

	var mtype string
	var mmode string
	var value interface{}

	// Set defaults from prometheus metrics data.
	switch {
	case dtoMetric.Counter != nil:
		mtype = "float"
		mmode = "counter"
		value = dtoMetric.Counter.GetValue()
	case dtoMetric.Gauge != nil:
		mtype = "float"
		mmode = "gauge"
		value = dtoMetric.Gauge.GetValue()
	default:
		return nil
	}

	// Apply overrides from alias mapping if set.
	if pm.mtype != "" {
		mtype = pm.mtype
	}
	if pm.mmode != "" {
		mmode = pm.mmode
	}

	m, err := ksurveyclient.NewConstMap(pm.fqName, map[string]interface{}{
		"desc":  pm.help,
		"mode":  mmode,
		"type":  mtype,
		"value": value,
	})
	if err != nil {
		return err
	}

	return m.Write(md)
}

type registry struct {
	registry  *ksurveyclient.Registry
	whitelist map[string]string
}

// WrapRegistry wraps the provided ksurveyclient.Registry so it can register
// prometheus.Collectors with optional filter and aliasing with the provided
// whitelist map.
//
// Example whitelist with formatting:
//
// map[string]string{
//	"rtm_distinct_users_connected_max": "usercnt_active,gauge:int",
//	"rtm_group_channels_created_max":   "usercnt_room,gauge:int",
//	"rtm_channels_created_max":         "usercnt_equipment,gauge:int",
//	"rtm_connections_connected_max":    "usercnt_nonactive,gauge:int",
//	"guest_http_logon_success_total":   "usercnt_na_user,gauge:int",
// }
//
func WrapRegistry(reg *ksurveyclient.Registry, whitelist map[string]string) prometheus.Registerer {
	if reg == nil {
		reg = DefaultRegistry
	}

	return &registry{
		registry:  reg,
		whitelist: whitelist,
	}
}

// Register registers the provided Collector with the associated Registry.
func (reg *registry) Register(c prometheus.Collector) error {
	return reg.registry.Register(WrapCollector(c, reg.whitelist))
}

// MustRegister registers the provided Collectors with the accociated Registry
// and panics if any error occurs.
func (reg *registry) MustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		if err := reg.Register(c); err != nil {
			panic(err)
		}
	}
}

// Unregister is there to satisfy the prometheus.Registerer interface. It has
// no effect on this registry.
func (reg *registry) Unregister(c prometheus.Collector) bool {
	return false
}
