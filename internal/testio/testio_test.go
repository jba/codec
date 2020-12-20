// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testio

import (
	"context"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"testing"
	"time"
)

const wantTput = 100

func TestThroughputReader(t *testing.T) {
	tr := NewThroughputReader(context.Background(), reader(150*meb), wantTput)
	checkThroughput(t, ioutil.Discard, tr, wantTput, nil)
}

func TestThroughputReaderSmallBurst(t *testing.T) {
	tr := NewThroughputReader(context.Background(), reader(150*meb), wantTput)
	checkThroughput(t, ioutil.Discard, tr, wantTput, make([]byte, 2*burst+17))
}

func TestThroughputWriter(t *testing.T) {
	tw := NewThroughputWriter(context.Background(), ioutil.Discard, wantTput)
	checkThroughput(t, tw, reader(150*meb), wantTput, nil)
}

func TestThroughputWriterSmallBurst(t *testing.T) {
	tw := NewThroughputWriter(context.Background(), ioutil.Discard, wantTput)
	checkThroughput(t, tw, reader(150*meb), wantTput, make([]byte, 2*burst+17))
}

func reader(size int64) io.Reader {
	return &io.LimitedReader{R: rand.New(rand.NewSource(0)), N: size}
}

func checkThroughput(t *testing.T, w io.Writer, r io.Reader, want int64, buf []byte) {
	t.Helper()
	start := time.Now()
	n, err := io.CopyBuffer(w, r, buf)
	if err != nil {
		t.Fatal(err)
	}
	got := float64(n) / meb / time.Since(start).Seconds()
	if math.Abs(float64(got-float64(want))) > 1 {
		t.Errorf("got throughput %f, want %d M/s", got, want)
	}

}
