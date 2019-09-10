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

func TestHashGUIDv1(t *testing.T) {
	raw := []byte("test-guid-1")
	guid := HashGUIDv1(raw)
	if string(guid) != "36D2CF432EA89B034D145311D7B001BF4AF425666EEB7652896A6C5A54975D16" {
		t.Fatalf("unexpected value: %v", string(guid))
	}
}
