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

package counter

import (
	"testing"
)

func TestUint(t *testing.T) {
	c := GetUint()
	if v := c.Value(); v != 0 {
		t.Fatalf("initial value %v is not 0", v)
	}
	c.Inc()
	if v := c.Value(); v != 1 {
		t.Errorf("value after first inc %v is not 1", v)
	}
	c.Inc()
	if v := c.Value(); v != 2 {
		t.Errorf("value after second inc %v is not 2", v)
	}
	c.Add(8)
	if v := c.Value(); v != 10 {
		t.Errorf("value after add %v is not the expected value", v)
	}
	c.Set(42)
	if v := c.Value(); v != 42 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
	c.Dec()
	if v := c.Value(); v != 41 {
		t.Errorf("value after dec %v is not the expected value", v)
	}
	c.Set(0)
	c.Dec()
	if v := c.Value(); v != 0 {
		t.Errorf("value after dec %v from 0 is not 0", v)
	}
}

func TestUintMax(t *testing.T) {
	c := GetUintMax()
	if v := c.Value(); v != 0 {
		t.Fatalf("initial value %v is not 0", v)
	}
	c.Set(42)
	if v := c.Value(); v != 42 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
	c.Set(100)
	if v := c.Value(); v != 100 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
	c.Set(50)
	if v := c.Value(); v != 100 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
}

func TestUintMin(t *testing.T) {
	c := GetUintMin()
	if v := c.Value(); v != maxUint64 {
		t.Fatalf("initial value %v is not max", v)
	}
	c.Set(42)
	if v := c.Value(); v != 42 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
	c.Set(100)
	if v := c.Value(); v != 42 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
	c.Set(5)
	if v := c.Value(); v != 5 {
		t.Fatalf("value after set %v is not the expected value", v)
	}
}
