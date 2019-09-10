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

// ConstMap is a used to hold a constant set of key values.
type ConstMap interface {
	Metric
	Collector
}

type constMap struct {
	name   string
	fields map[string]interface{}

	selfCollector
}

// NewConstMap creates a ConstMap with the provided name and field data.
func NewConstMap(name string, fields map[string]interface{}) (ConstMap, error) {
	// TODO(longsleep): Validate name and fields.
	cm := &constMap{
		name:   name,
		fields: fields,
	}
	cm.init(cm)
	return cm, nil
}

// MustNewConstMap creates a ConstMap with the provided name and field data
// and panics of an error occurs.
func MustNewConstMap(name string, fields map[string]interface{}) ConstMap {
	cm, err := NewConstMap(name, fields)
	if err != nil {
		panic(err)
	}
	return cm
}

func (mdm *constMap) Write(md *MetricData) error {
	md.Name = mdm.name
	md.Fields = mdm.fields

	return nil
}
