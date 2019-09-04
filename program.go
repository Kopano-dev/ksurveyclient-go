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
	"os"
	"path"
)

// ProgramVersion defines the default proram version number that is being used
// when not provided directly. It is exposed here so it can be set at compile
// time.
var ProgramVersion = "0.0.0-unknown"

type programCollector struct {
	name    string
	version string
}

// NewProgramCollector creates a Collector which collects program information.
func NewProgramCollector(name string, version string) Collector {
	if name == "" {
		name = path.Base(os.Args[0])
	}
	if version == "" {
		version = ProgramVersion
	}

	return &programCollector{
		name:    name,
		version: version,
	}
}

// Collect first gathers the associated managers collectors managers data. Then
// it creates constant metrics based on the returned data.
func (pc *programCollector) Collect(ch chan<- Metric) {
	ch <- MustNewConstMapMetric("program_name", map[string]interface{}{
		"desc":  "Program name",
		"type":  "string",
		"value": pc.name,
	})

	ch <- MustNewConstMapMetric("program_version", map[string]interface{}{
		"desc":  "Program name",
		"type":  "string",
		"value": pc.version,
	})
}
