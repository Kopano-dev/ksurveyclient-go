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
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	DefaultConfig.StartDelay = 1
	os.Exit(m.Run())
}

func TestStartKSurveyClient(t *testing.T) {
	ok := false
	ts := httptest.NewTLSServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("User-Agent") != DefaultConfig.UserAgent {
			t.Errorf("unexpected User-Agent: %v in request", req.Header.Get("User-Agent"))
		}
		// TODO(longsleep): Validate incoming request data.
		ok = true
	}))
	defer ts.Close()

	config := DefaultConfig.Clone()
	config.URL = ts.URL
	config.HTTPClient = ts.Client()
	config.Logger = &testingLogger{t}

	err := StartKSurveyClient(context.Background(), config, nil)
	if err != nil {
		t.Error("failed to start survey client", err)
	}

	select {
	case <-time.After(5 * time.Second):
	}
	if !ok {
		t.Error("request was not received")
	}
}
