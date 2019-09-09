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

package autosurvey

import (
	"context"
	"errors"
	"os"
	"sync"

	"stash.kopano.io/kgol/ksurveyclient-go"
)

var disabled = false
var mutex sync.Mutex

// DefaultRegistry exposes the registry which is used by autosurvey.
var DefaultRegistry *ksurveyclient.Registry

func init() {
	DefaultRegistry = ksurveyclient.DefaultRegistry

	if v := os.Getenv("KOPANO_SURVEYCLIENT_AUTOSURVEY"); v == "false" || v == "no" {
		disabled = true
		return
	}
}

var started = false

// Start is the function which gets auto survey up and running using the default
// configuration and the default registry with some standard collectors.
func Start(ctx context.Context, name, version string, cs ...ksurveyclient.Collector) error {
	return start(ctx, name, version, cs...)
}

// MustStart is the function which gets auto survey with Start up and running
// but panics if start fails.
func MustStart(ctx context.Context, name, version string, cs ...ksurveyclient.Collector) {
	err := Start(ctx, name, version, cs...)
	if err != nil {
		panic(err)
	}
}

func start(ctx context.Context, name, version string, cs ...ksurveyclient.Collector) error {
	mutex.Lock()
	defer mutex.Unlock()
	if started {
		return errors.New("already started")
	}
	started = true
	if disabled {
		return nil
	}

	reg := DefaultRegistry
	err := reg.Register(ksurveyclient.NewProgramCollector(name, version))
	if err != nil {
		return nil
	}
	for _, c := range cs {
		err = reg.Register(c)
		if err != nil {
			return err
		}
	}
	return ksurveyclient.StartKSurveyClient(ctx, nil, nil)
}
