// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	"github.com/jba/codec/internal/bench"
	"github.com/jba/codec/internal/testio"
	ucodec "github.com/ugorji/go/codec"
)

//go:generate rm -f types.gen.go
//go:generate go run . gen code

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
	jbaCodec = Codec{
		"jba/codec orig",
		func(w io.Writer, data interface{}) error {
			e := codecapi.NewEncoder(w, codecapi.EncodeOptions{})
			return e.Encode(data)
		},
		jbaCodecDecode,
	}
)

var codecs = []Codec{
	jbaCodec,
	{
		"gob",
		func(w io.Writer, data interface{}) error {
			e := gob.NewEncoder(w)
			return e.Encode(data)
		},
		func(r io.Reader, ptr interface{}) error {
			d := gob.NewDecoder(r)
			return d.Decode(ptr)
		},
	},
	{
		"ugorji-cbor",
		func(w io.Writer, data interface{}) error {
			e := ucodec.NewEncoder(w, &ucodec.CborHandle{})
			return e.Encode(data)
		},
		func(r io.Reader, ptr interface{}) error {
			return ucodec.NewDecoder(r, &ucodec.CborHandle{}).Decode(ptr)
		},
	},
	// ugorji with msgpack and binc have almost identical performance to ugorji with cbor
	// {
	// 	"ugorji-msgpack",
	// 	func(w io.Writer, data interface{}) error {
	// 		e := ucodec.NewEncoder(w, &ucodec.MsgpackHandle{})
	// 		return e.Encode(data)
	// 	},
	// 	func(r io.Reader, ptr interface{}) error {
	// 		return ucodec.NewDecoder(r, &ucodec.MsgpackHandle{}).Decode(ptr)
	// 	},
	// },
	// {
	// 	"ugorji-binc",
	// 	func(w io.Writer, data interface{}) error {
	// 		e := ucodec.NewEncoder(w, &ucodec.BincHandle{})
	// 		return e.Encode(data)
	// 	},
	// 	func(r io.Reader, ptr interface{}) error {
	// 		return ucodec.NewDecoder(r, &ucodec.BincHandle{}).Decode(ptr)
	// 	},
	// },
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

type benchmarkData struct {
	name   string
	read   func() (interface{}, error)
	newptr func() interface{}
}

var datas = []benchmarkData{
	{
		"hyperledger",
		func() (interface{}, error) { return hlDecodeJSON("ledgerAPIs.json") },
		func() interface{} { var sb submittedData; return &sb },
	},
	{
		"licenses",
		func() (interface{}, error) { var ld LicenseData; return gobDecodeFile("licenses.gob", &ld) },
		func() interface{} { var ld *LicenseData; return &ld },
	},
	{
		"licenses-small",
		func() (interface{}, error) { var ld LicenseData; return gobDecodeFile("licenses-small.gob", &ld) },
		func() interface{} { var ld *LicenseData; return &ld },
	},
	{
		"stocks",
		func() (interface{}, error) {
			var sds []*StockData
			if _, err := gobDecodeFile("stocks.gob", &sds); err != nil {
				return nil, err
			}
			return sds, nil
		},
		func() interface{} { var sds []*StockData; return &sds },
	},
	{
		"scores",
		func() (interface{}, error) {
			var sds []Score
			if _, err := gobDecodeFile("scores.gob", &sds); err != nil {
				return nil, err
			}
			return sds, nil
		},
		func() interface{} { return new([]Score) },
	},
}

var cpuProfileFile *os.File

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "gen":
		things := flag.Args()[1:]
		if len(things) == 0 {
			log.Fatal("need things to generate")
		}
		if err := generate(things); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Generated files, exiting.")

	case "bet":
		runBreakEvenThroughput()

	case "bm":
		runBenchmarks()

	default:
		log.Fatalf("unknown command %q", flag.Arg(0))
	}
}

func runBenchmarks() {
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

	for _, bd := range datas {
		runBenchmark(bd)
	}

	if *allocprofile {
		kind := "allocs"
		hp := pprof.Lookup(kind)
		err := writeNewFile(kind+".out", func(f *os.File) error { return hp.WriteTo(f, 0) })
		if err != nil {
			log.Fatal(err)
		}
	}
}

// runBenchmark uses bd to read data to be used for benchmarks.
// It then uses the data to measure encoding and decoding for all the codecs.
func runBenchmark(bd benchmarkData) {
	data, err := bd.read()
	if err != nil {
		log.Fatal(err)
	}
	for _, tput := range throughputs {
		s := "max"
		if tput > 0 {
			s = fmt.Sprintf("%d", tput)
		}
		fmt.Printf("---- %s at %s Mi/sec ----\n", bd.name, s)
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
			bms = append(bms, newDecodeBenchmark(encoded[i], c, tput, bd.newptr))
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

func gobDecodeFile(filename string, ptr interface{}) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := gob.NewDecoder(f)
	if err := d.Decode(ptr); err != nil {
		return nil, err
	}
	return ptr, nil
}

func totalAlloc() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.TotalAlloc
}

func writeNewFile(filename string, writer func(*os.File) error) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := writer(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
