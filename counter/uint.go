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
	"errors"
	"sync/atomic"
)

// Uint is an interface for positive integer counters.
type Uint interface {
	// Inc increments the counter by 1.
	Inc()
	// Dec increments the counter by 1.
	Dec()
	// Add increments the counter by the provided value.
	Add(uint64)
	// Set sets the counters value.
	Set(uint64)
	// Value returns the counters current value.
	Value() uint64
}

// UintMinMax is an interface for minimal and maximal positive integer counters.
type UintMinMax interface {
	// Set sets the counters value.
	Set(uint64)
	// Value returns the counters current value.
	Value() uint64
}

type uintImpl struct {
	value uint64
}

func (c *uintImpl) Inc() {
	atomic.AddUint64(&c.value, 1)
}

func (c *uintImpl) Dec() {
	for {
		if o := atomic.LoadUint64(&c.value); o > 0 {
			if !atomic.CompareAndSwapUint64(&c.value, o, o-1) {
				continue
			}
		}
		break
	}
}

func (c *uintImpl) Add(v uint64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in value"))
	}
	atomic.AddUint64(&c.value, v)
}

func (c *uintImpl) Set(num uint64) {
	atomic.StoreUint64(&c.value, num)
}

func (c *uintImpl) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// GetUint returns a new counter to count positive integers.
func GetUint() Uint {
	return &uintImpl{}
}

type uintMaxImpl uintImpl

func (c *uintMaxImpl) Set(v uint64) {
	for {
		if o := atomic.LoadUint64(&c.value); v > o {
			if !atomic.CompareAndSwapUint64(&c.value, o, v) {
				continue
			}
		}
		break
	}
}

func (c *uintMaxImpl) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// GetUintMax returns a new counter to count maximum positive integers.
func GetUintMax() UintMinMax {
	return &uintMaxImpl{}
}

type uintMinImpl uintImpl

func (c *uintMinImpl) Set(v uint64) {
	for {
		if o := atomic.LoadUint64(&c.value); v < o {
			if !atomic.CompareAndSwapUint64(&c.value, o, v) {
				continue
			}
		}
		break
	}
}

func (c *uintMinImpl) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

const maxUint64 = ^uint64(0)

// GetUintMin returns a new counter to count minimum positive integers.
func GetUintMin() UintMinMax {
	return &uintMinImpl{value: maxUint64}
}
