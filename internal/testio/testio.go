// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testio provides Readers and Writers for testing.
package testio

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

const meb = 1024 * 1024

// A ThroughputReader limits read throughput.
type ThroughputReader struct {
	ctx context.Context
	r   io.Reader
	lim *rate.Limiter
}

// NewThroughputReader creates a ThroughputReader that limits throughput to
// msec mebibytes per second.
func NewThroughputReader(ctx context.Context, r io.Reader, msec int) *ThroughputReader {
	return &ThroughputReader{ctx, r, limiter(msec)}
}

// Create a limiter that uses bytes/sec, with a burst of 1M.
func limiter(msec int) *rate.Limiter {
	limit := rate.Inf
	if msec > 0 {
		limit = rate.Limit(msec * meb)
	}
	return rate.NewLimiter(limit, meb)
}

// Read implements io.Reader.
func (r *ThroughputReader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if err := r.lim.WaitN(r.ctx, n); err != nil {
		return n, err
	}
	return n, err
}

// A ThroughputWriter limits write throughput.
type ThroughputWriter struct {
	ctx context.Context
	w   io.Writer
	lim *rate.Limiter
}

// NewThroughputWriter creates a ThroughputWriter that limits throughput to
// msec mebibytes per second.
func NewThroughputWriter(ctx context.Context, w io.Writer, msec int) *ThroughputWriter {
	return &ThroughputWriter{ctx, w, limiter(msec)}
}

// Write implements io.Write.
func (w *ThroughputWriter) Write(buf []byte) (int, error) {
	b := w.lim.Burst()
	n := 0
	for len(buf) > 0 {
		var wbuf []byte
		if len(buf) < b {
			wbuf, buf = buf, nil
		} else {
			wbuf, buf = buf[:b], buf[b:]
		}
		if err := w.lim.WaitN(w.ctx, len(wbuf)); err != nil {
			return n, err
		}
		m, err := w.w.Write(wbuf)
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
