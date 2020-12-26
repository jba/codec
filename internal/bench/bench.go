// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bench is for comparing benchmarks.
package bench

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"text/tabwriter"
	"time"
)

// A Benchmark is a named benchmarking function.
type Benchmark struct {
	Name string
	Func func(b *testing.B) error
}

// Run runs the given benchmarks and write a report of the results. The speed of
// each benchmark after the first is displayed as a multiplier of the first's
// speed.
func Run(bms []Benchmark) {
	var r0 testing.BenchmarkResult
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 2, ' ', tabwriter.AlignRight)
	// Do a warm-up run.
	for _, bm := range bms {
		_, _ = Run1(bm)
	}
	for i, bm := range bms {
		r, err := Run1(bm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", bm.Name, err)
		}
		if i == 0 {
			r0 = r
		}
		if err == nil {
			d := time.Duration(r.NsPerOp())
			fmt.Fprintf(w, "%s\t%d\t%dK/op\t%.2fs/op\t %.2fx\n",
				bm.Name, r.N, r.AllocedBytesPerOp()/1024, d.Seconds(), float64(r0.NsPerOp())/float64(r.NsPerOp()))
		}
	}
	w.Flush()
}

// Run1 runs bm.
func Run1(bm Benchmark) (testing.BenchmarkResult, error) {
	var err error
	runtime.GC()
	r := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		err = bm.Func(b)
		if err != nil {
			b.Fatal(err)
		}
	})
	return r, err
}
