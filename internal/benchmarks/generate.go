// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"

	"github.com/jba/codec"
)

func generate(words []string) error {
	for _, w := range words {
		var err error
		switch w {
		case "code":
			err = codec.GenerateFile("types.gen.go", "main",
				LicenseData{}, submittedData{}, []*StockData(nil), []Score(nil))
			if err == nil {
				err = codec.GenerateFile("ast_types.gen.go", "main", append(astTypes, map[string]*ast.File{})...)
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
