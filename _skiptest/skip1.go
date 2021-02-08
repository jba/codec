// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build version1

package main

type Skip struct {
	U  uint64
	S1 string
	S2 string
	L  []*Skip
}
