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

package ksurveyclient // import "stash.kopano.io/kgol/ksurveyclient-go"

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SurveyClientEnabled directly controls if surveys are sent. When false, all
// survice clients do nothing.
var SurveyClientEnabled = true

func init() {
	if v := os.Getenv("KOPANO_SURVEYCLIENT_ENABLED"); v == "false" || v == "no" {
		SurveyClientEnabled = false
	}
}

type kSurveyClient struct {
	url        *url.URL
	startDelay uint64
	errorDelay uint64
	interval   uint64
	userAgent  string

	registry *Registry

	client *http.Client
	logger logger
}

// StartKSurveyClient starts a new survey client using the provided Context and
// the provid Config.
func StartKSurveyClient(ctx context.Context, config *Config, registry *Registry) error {
	var err error

	if config == nil {
		config = DefaultConfig
	}
	if registry == nil {
		registry = DefaultRegistry
	}

	ksv := &kSurveyClient{
		startDelay: config.StartDelay,
		interval:   config.Interval,
		userAgent:  config.UserAgent,

		registry: registry,

		logger: config.Logger,
	}
	if ksv.logger == nil {
		ksv.logger = DefaultLogger
	}
	if ksv.url, err = url.Parse(config.URL); err != nil {
		return err
	}

	if config.HTTPClient != nil {
		if config.Insecure {
			return errors.New("inconsistent configuration, either set HTTPClient or Insecure")
		}
		ksv.client = config.HTTPClient
	} else {
		var tlsClientConfig *tls.Config
		if config.Insecure {
			tlsClientConfig = &tls.Config{
				InsecureSkipVerify: config.Insecure,
			}
		}
		ksv.client = &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: -1,
					DualStack: true,
				}).DialContext,
				DisableKeepAlives:     true,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       tlsClientConfig,
			},
		}
	}

	go ksv.Run(ctx)

	return nil
}

func (ksv *kSurveyClient) Run(ctx context.Context) {
	if ksv.startDelay > 0 {
		select {
		case <-ctx.Done():
			// Context done, exit.
			return
		case <-time.After(time.Duration(ksv.startDelay) * time.Second):
			// Continue after start delay.
		}
	}
	var err error
	var interval uint64
	for {
		interval = ksv.interval
		err = ksv.Do()
		if err != nil && ksv.logger != nil {
			ksv.logger.Printf("ksurveyclient failed: %v", err)
			if ksv.errorDelay > 0 {
				interval = ksv.errorDelay
			}
		}
		if interval == 0 {
			// Done.
			return
		}
		select {
		case <-ctx.Done():
			// Context done, exit.
			return
		case <-time.After(time.Duration(ksv.interval) * time.Second):
			// Continue after interval.
		}
	}
}

func (ksv *kSurveyClient) Do() error {
	if !SurveyClientEnabled {
		// Global disable flag - do nothing.
		return nil
	}

	ms, err := ksv.registry.Gather()
	if err != nil {
		return err
	}
	payload := kSurveyPayloadV2{
		Version: 2,
		Stats:   ms,
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ksv.url.String(), &buf)
	req.Header.Set("X-Kopano-Stats-Request", "1")
	req.Header.Set("Content-Type", "application/json")
	if ksv.userAgent != "" {
		req.Header.Set("User-Agent", ksv.userAgent)
	}

	resp, err := ksv.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return err
}
