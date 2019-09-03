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
	"net/http"
	"os"
	"strconv"
)

// Config defines the settings for the service client.
type Config struct {
	URL        string
	StartDelay uint64
	ErrorDelay uint64
	Interval   uint64
	Insecure   bool

	Logger     logger
	HTTPClient *http.Client
}

// Clone returns a copy of the associated Config.
func (c *Config) Clone() *Config {
	return &Config{
		URL:        c.URL,
		StartDelay: c.StartDelay,
		ErrorDelay: c.ErrorDelay,
		Interval:   c.Interval,
		Insecure:   c.Insecure,

		Logger: c.Logger,
	}
}

// DefaultConfig hols the service client default configuration.
var DefaultConfig = &Config{
	URL:        "https://stats.kopano.io/api/stats/v1/submit",
	StartDelay: 60,
	ErrorDelay: 60,
	Interval:   3600,
	Insecure:   false,
}

func init() {
	if v := os.Getenv("KOPANO_SURVEYCLIENT_URL"); v != "" {
		DefaultConfig.URL = v
	}
	if v := os.Getenv("KOPANO_SURVEYCLIENT_START_DELAY"); v != "" {
		DefaultConfig.StartDelay, _ = strconv.ParseUint(v, 10, 64)
	}
	if v := os.Getenv("KOPANO_SURVEYCLIENT_ERROR_DELAY"); v != "" {
		DefaultConfig.ErrorDelay, _ = strconv.ParseUint(v, 10, 64)
	}
	if v := os.Getenv("KOPANO_SURVEYCLIENT_INTERVAL"); v != "" {
		DefaultConfig.Interval, _ = strconv.ParseUint(v, 10, 64)
	}
	if v := os.Getenv("KOPANO_SURVEYCLIENT_INSECURE"); v != "" {
		DefaultConfig.Insecure = v == "yes"
	}
}
