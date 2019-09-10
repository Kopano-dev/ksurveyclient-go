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
	"bytes"
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

// HashGUIDv1 takes the provided GUID and returns its hash.
func HashGUIDv1(rawGUID []byte) []byte {
	value := blake2b.Sum256(rawGUID)
	out := make([]byte, hex.EncodedLen(len(value)))
	hex.Encode(out, value[:])
	return bytes.ToUpper(out)
}
