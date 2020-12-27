// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/* TODO

Maybe generate more stuff for the hyperledger benchmark. A lot of types and
fields are unused?


Other possible benchmarks:
- https://github.com/robertkrimen/otto/blob/15f95af6e78dcd2030d8195a138bd88d4f403546/script.go
- https://github.com/gocolly/colly/blob/1cd684083cf9bf9a8e33b5dfd6414d8516ae63af/http_backend.go#L161
*/

// A program for benchmarking codecs.
package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"testing"

	"github.com/jba/codec/codecapi"
	"github.com/jba/codec/internal/benchmarks/bench"
	"github.com/jba/codec/internal/benchmarks/data"
	"github.com/jba/codec/internal/benchmarks/testio"
	ucodec "github.com/ugorji/go/codec"
)

var _ ucodec.CborHandle

var (
	cpuprofile   = flag.String("cpuprofile", "", "write cpu profile to `file`")
	allocprofile = flag.Bool("allocprofile", false, "write alloc profile")
)

// Throughputs to benchmark, in Mi/sec.
var throughputs = []int{
	0,    // unlimited throughput; speed of memory
	3000, // reading from local disk
	250,  // reading from a GCS bucket
	100,  // reading from a cloud DB
}

type Codec struct {
	name   string
	encode func(io.Writer, interface{}) error
	decode func(io.Reader, interface{}) error
}

var (
	jbaCodec    = newJbaCodec("", codecapi.EncodeOptions{})
	jbaCodecBig = newJbaCodec("big", codecapi.EncodeOptions{Buffer: make([]byte, 100*1024*1024)})
	gobCodec    = Codec{
		"gob",
		func(w io.Writer, data interface{}) error {
			e := gob.NewEncoder(w)
			return e.Encode(data)
		},
		func(r io.Reader, ptr interface{}) error {
			d := gob.NewDecoder(r)
			return d.Decode(ptr)
		},
	}
	ugorjiCodec = Codec{
		"ugorji-cbor",
		func(w io.Writer, data interface{}) error {
			e := ucodec.NewEncoder(w, &ucodec.CborHandle{})
			return e.Encode(data)
		},
		func(r io.Reader, ptr interface{}) error {
			return ucodec.NewDecoder(r, &ucodec.CborHandle{}).Decode(ptr)
		},
	}
)

var codecs = []Codec{
	jbaCodec,
	jbaCodecBig,
	gobCodec,
	ugorjiCodec,
	// ugorji with msgpack and binc have almost identical performance to ugorji with cbor.
}

func newJbaCodec(suffix string, opts codecapi.EncodeOptions) Codec {
	name := "jba/codec"
	if suffix != "" {
		name += " " + suffix
	}
	return Codec{
		name,
		func(w io.Writer, data interface{}) error {
			e := codecapi.NewEncoder(w, opts)
			return e.Encode(data)
		},
		jbaCodecDecode,
	}
}

func jbaCodecDecode(r io.Reader, ptr interface{}) error {
	d := codecapi.NewDecoder(r)
	x, err := d.Decode()
	if err != nil {
		return err
	}
	reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(x))
	return nil
}

var datas = []data.BenchmarkData{
	data.ASTData,
	data.Hyperledger,
	data.Licenses,
	data.LicensesSmall,
	data.Stocks,
	data.Scores,
}

var cpuProfileFile *os.File

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "bet":
		runBreakEvenThroughput(flag.Args()[1:])

	case "bm":
		runBenchmarks(flag.Args()[1:])

	default:
		log.Fatalf("unknown command %q", flag.Arg(0))
	}
}

func runBenchmarks(dataNames []string) {
	if *cpuprofile != "" {
		var err error
		cpuProfileFile, err = os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer func() {
			if err := cpuProfileFile.Close(); err != nil {
				log.Fatal("closing CPU profile file: ", err)
			}
		}()
	}

	for _, bd := range datasToRun(dataNames) {
		runBenchmark(bd)
	}

	if *allocprofile {
		kind := "allocs"
		hp := pprof.Lookup(kind)
		err := data.WriteNewFile(kind+".out", func(f *os.File) error { return hp.WriteTo(f, 0) })
		if err != nil {
			log.Fatal(err)
		}
	}
}

func datasToRun(dataNames []string) []data.BenchmarkData {
	if len(dataNames) == 0 {
		return datas
	}
	runName := map[string]bool{}
	for _, n := range dataNames {
		runName[n] = true
	}
	var ds []data.BenchmarkData
	for _, bd := range datas {
		if runName[bd.Name] {
			ds = append(ds, bd)
		}
	}
	return ds
}

// runBenchmark uses bd to read data to be used for benchmarks.
// It then uses the data to measure encoding and decoding for all the codecs.
func runBenchmark(bd data.BenchmarkData) {
	data, err := bd.Read()
	if err != nil {
		log.Fatalf("%s: %v", bd.Name, err)
	}
	for _, tput := range throughputs {
		s := "max"
		if tput > 0 {
			s = fmt.Sprintf("%d", tput)
		}
		fmt.Printf("---- %s at %s Mi/sec ----\n", bd.Name, s)
		fmt.Println("encode")
		var bms []bench.Benchmark
		encoded := make([][]byte, len(codecs))
		for p, c := range codecs {
			bms = append(bms, newEncodeBenchmark(data, c, tput, &encoded[p]))
		}
		bench.Run(bms)
		fmt.Println()

		fmt.Println("decode")
		bms = nil
		for i, c := range codecs {
			c := c
			bms = append(bms, newDecodeBenchmark(encoded[i], c, tput, bd.Newptr))
		}
		bench.Run(bms)
		fmt.Println()
	}
}

func newEncodeBenchmark(data interface{}, c Codec, tput int, out *[]byte) bench.Benchmark {
	return bench.Benchmark{
		Name: c.name,
		Func: func(b *testing.B) error {
			var buf bytes.Buffer
			w := testio.NewThroughputWriter(context.Background(), &buf, tput)
			for i := 0; i < b.N; i++ {
				buf.Reset()
				if err := c.encode(w, data); err != nil {
					return err
				}
			}
			*out = buf.Bytes()
			return nil
		},
	}
}

func newDecodeBenchmark(enc []byte, c Codec, tput int, newptr func() interface{}) bench.Benchmark {
	return bench.Benchmark{
		Name: c.name,
		Func: func(b *testing.B) error {
			br := bytes.NewReader(enc)
			r := testio.NewThroughputReader(context.Background(), br, tput)
			for i := 0; i < b.N; i++ {
				br.Reset(enc)
				profiling := false
				if i == 0 && c.name == "codec" && cpuProfileFile != nil {
					fmt.Println("starting cpu profile")
					profiling = true
					if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
						log.Fatal(err)
					}
				}
				if err := c.decode(r, newptr()); err != nil {
					return err
				}
				if profiling {
					pprof.StopCPUProfile()
					os.Exit(0)
				}
			}
			return nil
		},
	}
}

func totalAlloc() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.TotalAlloc
}
