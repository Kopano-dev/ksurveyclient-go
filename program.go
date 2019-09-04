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

// Default proram name and version number which are being used when not provided
// directly or empty. It is exposed here so it can be overriden.
var (
	DefaultProgramName    = "unknown"
	DefaultProgramVersion = "0.0.0-unknown"
)

type programCollector struct {
	name    string
	version string
}

// NewProgramCollector creates a Collector which collects program information.
func NewProgramCollector(name string, version string) Collector {
	if name == "" {
		name = path.Base(os.Args[0])
	}

	return &programCollector{
		name:    name,
		version: version,
	}
}

// Collect first gathers the associated managers collectors managers data. Then
// it creates constant metrics based on the returned data.
func (pc *programCollector) Collect(ch chan<- Metric) {
	func() {
		name := pc.name
		if pc.name == "" {
			name = DefaultProgramName
		}
		ch <- MustNewConstMapMetric("program_name", map[string]interface{}{
			"desc":  "Program name",
			"type":  "string",
			"value": name,
		})
	}()

	func() {
		version := pc.version
		if version == "" {
			version = DefaultProgramVersion
		}
		ch <- MustNewConstMapMetric("program_version", map[string]interface{}{
			"desc":  "Program version",
			"type":  "string",
			"value": version,
		})
	}()
}
