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
	"os"

	"stash.kopano.io/kgol/ksurveyclient-go"
)

func init() {
	if v := os.Getenv("KOPANO_SURVEYCLIENT_AUTOSURVEY"); v == "false" || v == "no" {
		return
	}

	reg := ksurveyclient.DefaultRegistry
	reg.MustRegister(ksurveyclient.NewProgramCollector("", ""))

	ksurveyclient.StartKSurveyClient(context.Background(), nil, nil)
}

// SetProgramNameAndVersion allows to sets the program collector data when using
// autosurvey.
func SetProgramNameAndVersion(name, version string) {
	if name != "" {
		ksurveyclient.DefaultProgramName = name
	}
	if version != "" {
		ksurveyclient.DefaultProgramVersion = version
	}
}
