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
	"testing"
)

var testGUID = []byte("test-guid")

func TestProgramCollector(t *testing.T) {
	pc := NewProgramCollector("", "", testGUID)
	var metricChan = make(chan Metric)
	go func() {
		pc.Collect(metricChan)
		close(metricChan)
	}()
	for metric := range metricChan {
		md := &MetricData{}
		metric.Write(md)
		switch md.Name {
		case "program_name":
			value := md.Fields["value"].(string)
			if value == "" {
				t.Error("program_name is empty")
			}
		case "program_version":
			value := md.Fields["value"].(string)
			if value != "0.0.0-unknown" {
				t.Errorf("unexpected program_version: %v", value)
			}
		case "server_guid":
			value := md.Fields["value"].(string)
			if value != string(testGUID) {
				t.Errorf("unexpected server_guid: %v", value)
			}
		default:
			t.Errorf("unknown program metrics name: %v", md.Name)
		}
	}
}
