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

package ksurveyclient

import (
	"sync"
)

// A Registry holds registered Collectors and collects their Metrics.
type Registry struct {
	collectors []Collector
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		collectors: make([]Collector, 0),
	}
}

// DefaultRegistry is a Registry which is used by default.
var DefaultRegistry *Registry

func init() {
	DefaultRegistry = NewRegistry()
	DefaultRegistry.MustRegister(NewBasicCollector())
}

// Register registers the provided Collector with the associated Registry.
func (reg *Registry) Register(c Collector) error {
	reg.collectors = append(reg.collectors, c)

	return nil
}

// MustRegister registers the provided Collectors with the accociated Registry
// and panics if any error occurs.
func (reg *Registry) MustRegister(cs ...Collector) {
	for _, c := range cs {
		if err := reg.Register(c); err != nil {
			panic(err)
		}
	}
}

// Gather calls the Collect method of the registered Collectors and then
// gathers the collected metrics into a MetricSet.
func (reg *Registry) Gather() (*MetricSet, error) {
	var wg sync.WaitGroup
	var num = len(reg.collectors)
	var metricChan = make(chan Metric, num)
	collectors := make(chan Collector, num)

	for _, collector := range reg.collectors {
		collectors <- collector
	}

	wg.Add(num)
	collectWorker := func() {
		for {
			select {
			case collector := <-collectors:
				collector.Collect(metricChan)
			default:
				return
			}
			wg.Done()
		}
	}
	go collectWorker()

	go func() {
		wg.Wait()
		close(metricChan)
	}()

	defer func() {
		if metricChan != nil {
			for range metricChan {
			}
		}
	}()

	content := make([]*MetricData, 0)
	cc := metricChan

	for {
		select {
		case metric, ok := <-cc:
			if !ok {
				cc = nil
				break
			}
			md := &MetricData{}
			if err := metric.Write(md); err == nil {
				content = append(content, md)
			}
		}

		if cc == nil {
			break
		}
	}

	return &MetricSet{
		Content: content,
	}, nil
}
