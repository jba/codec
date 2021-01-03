// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"reflect"
	"testing"
)

func BenchmarkMapStore(b *testing.B) {
	// Demonstrate that storing a pointer into a map[interface{}]int is
	// significantly slower than a map[uintptr]int.
	p := &b
	b.Run("interface", func(b *testing.B) {
		m := map[interface{}]int{}
		for i := 0; i < b.N; i++ {
			m[p] = 1
		}
	})
	b.Run("uintptr", func(b *testing.B) {
		m := map[uintptr]int{}
		for i := 0; i < b.N; i++ {
			m[reflect.ValueOf(p).Pointer()] = 1
		}
	})
}
