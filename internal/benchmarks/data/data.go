// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

//go:generate ./generate.sh

import (
	"encoding/gob"
	"fmt"
	"go/ast"
	"os"
	"reflect"

	"github.com/jba/codec"
)

type BenchmarkData struct {
	Name   string
	Read   func() (interface{}, error)
	Newptr func() interface{}
}

func gobBenchmarkData(name string, newptr func() interface{}) BenchmarkData {
	return BenchmarkData{
		Name:   name,
		Newptr: newptr,
		Read: func() (interface{}, error) {
			ptr := newptr()
			if _, err := gobDecodeFile("data/"+name+".gob", ptr); err != nil {
				return nil, err
			}
			return reflect.ValueOf(ptr).Elem().Interface(), nil
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

func WriteNewFile(filename string, writer func(*os.File) error) error {
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

func Generate(words []string) error {
	thisPkgPath := "github.com/jba/codec/internal/benchmarks/data"
	for _, w := range words {
		var err error
		switch w {
		case "code":
			err = codec.GenerateFile("types.gen.go", thisPkgPath, nil,
				&LicenseData{}, submittedData{}, []*StockData(nil), []Score(nil))
			if err == nil {
				err = codec.GenerateFile("ast_types.gen.go", thisPkgPath, nil,
					append(astTypes, map[string]*ast.File{})...)
			}
		case "ast":
			err = generateASTToFile("ast.gob")
		case "stocks":
			err = generateStockDataToFile("stocks.gob")
		case "scores":
			err = generateScoreDataToFile("scores.gob")
		case "licenses":
			err = writeSmallLicenseFile()
		default:
			err = fmt.Errorf("unknown generate kind %q", w)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
