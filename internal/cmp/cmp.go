// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a package with the same name as github.com/google/go-cmp/cmp.
package cmp

type Option int

func FuncOptionSlice() interface{} {
	type Option int
	return []Option{}
}
