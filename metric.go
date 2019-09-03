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
	"encoding/json"
)

// A MetricSet holds the collected MetricData.
type MetricSet struct {
	Content []*MetricData
}

// MarshalJSON serializes the associated MetricSet collected data to JSON.
func (ms *MetricSet) MarshalJSON() ([]byte, error) {
	cs := make(map[string]interface{})
	for _, c := range ms.Content {
		cs[c.Name] = c.Fields
	}

	return json.Marshal(cs)
}

// MetricData holds the collected data with its name and fields.
type MetricData struct {
	Name   string
	Fields map[string]interface{}
}

type constMapMetric struct {
	name   string
	fields map[string]interface{}
}

// NewConstMapMetric creates a metric with the provided name and field data.
func NewConstMapMetric(name string, fields map[string]interface{}) (Metric, error) {
	// TODO(longsleep): Validate name and fields.
	return &constMapMetric{
		name:   name,
		fields: fields,
	}, nil
}

// MustNewConstMapMetric creates a metric with the provided name and field data
// and panics of an error occurs.
func MustNewConstMapMetric(name string, fields map[string]interface{}) Metric {
	mdm, err := NewConstMapMetric(name, fields)
	if err != nil {
		panic(err)
	}
	return mdm
}

func (mdm *constMapMetric) Write(md *MetricData) error {
	md.Name = mdm.name
	md.Fields = mdm.fields

	return nil
}

// Metric is the interface implemented by anything that can be used to provide
// survey Metrics.
type Metric interface {
	// Write encodes the Metric into MetricData.
	Write(*MetricData) error
}
