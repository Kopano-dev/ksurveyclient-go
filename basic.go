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
	"bufio"
	"os"
	"strings"
	"syscall"
)

type basicCollector struct {
}

// NewBasicCollector creates a new Collector which collects basic information.
func NewBasicCollector() Collector {
	return &basicCollector{}
}

// Collect first gathers the associated managers collectors managers data. Then
// it creates constant metrics based on the returned data.
func (mc *basicCollector) Collect(ch chan<- Metric) {
	func() {
		var machineID string
		if f, err := os.Open("/etc/machine-id"); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f) // Default split function is ScanLines.
			if scanner.Scan() {
				machineID = scanner.Text()
			}
		}
		ch <- MustNewConstMapMetric("machine_id", map[string]interface{}{
			"desc":  "",
			"type":  "string",
			"value": machineID,
		})
	}()

	func() {
		var prettyName string
		var buf syscall.Utsname
		if err := syscall.Uname(&buf); err == nil {
			prettyName = charsToString(buf.Sysname[:]) + " " + charsToString(buf.Machine[:]) + " " + charsToString(buf.Release[:])
		}
		ch <- MustNewConstMapMetric("utsname", map[string]interface{}{
			"desc":  "Pretty platform name",
			"type":  "string",
			"value": prettyName,
		})
	}()

	func() {
		var prettyName string
		if f, err := os.Open("/etc/os-release"); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					prettyNameParts := strings.SplitAfterN(line, "PRETTY_NAME=", 2)
					if len(prettyNameParts) == 2 {
						prettyName = strings.Trim(prettyNameParts[1], "\" \t\r\n")
						break
					}
				}
			}
		} else if f, err := os.Open("/etc/redhat-release"); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			if scanner.Scan() {
				prettyName = scanner.Text()
			}
		}
		ch <- MustNewConstMapMetric("osrelease", map[string]interface{}{
			"desc":  "Pretty operating system name",
			"type":  "string",
			"value": prettyName,
		})
	}()
}
