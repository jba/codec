// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"os"

	"github.com/google/licensecheck"
)

var (
	// To reconstruct licenses.gob,
	//     cat licenses-gob-* > licenses.gob
	Licenses      = licenseBenchmarkData("licenses")
	LicensesSmall = licenseBenchmarkData("licenses-small")
)

func licenseBenchmarkData(name string) BenchmarkData {
	return BenchmarkData{
		name,
		func() (interface{}, error) {
			var ld LicenseData
			d, err := gobDecodeFile("data/"+name+".gob", &ld)
			return d, err
		},
		func() interface{} { return new(*LicenseData) },
	}
}

type LicenseData struct {
	Files    []*LicenseFile
	Contents []*LicenseContents
}

// LicenseFile holds information about a license file.
type LicenseFile struct {
	Module   string
	Version  string
	FilePath string
	Contents int // index into LicenseData.Contents; not a pointer because gob does not dedup
}

func (f1 *LicenseFile) Equal(f2 *LicenseFile) bool {
	return *f1 == *f2
}

// LicenseContents hold the contents of a license file, and information derived from it.
type LicenseContents struct {
	Contents     []byte
	ContentsHash [sha256.Size]byte     // SHA256 of the contents, to dedup equal contents
	OldTypes     []string              // from the DB, stored in the gob file
	OldCoverage  licensecheck.Coverage // ditto
	NewTypes     []string              // not populated from the gob
	NewCoverage  licensecheck.Coverage // ditto
}

func (c1 *LicenseContents) Equal(c2 *LicenseContents) bool {
	return bytes.Equal(c1.Contents, c2.Contents) &&
		c1.ContentsHash == c2.ContentsHash &&
		stringsEqual(c1.OldTypes, c2.OldTypes) &&
		coverageEqual(c1.OldCoverage, c2.OldCoverage) &&
		stringsEqual(c1.NewTypes, c2.NewTypes) &&
		coverageEqual(c1.NewCoverage, c2.NewCoverage)
}

func coverageEqual(c1, c2 licensecheck.Coverage) bool {
	if c1.Percent != c2.Percent {
		return false
	}
	if len(c1.Match) != len(c2.Match) {
		return false
	}
	for i, m1 := range c1.Match {
		if m1 != c2.Match[i] {
			return false
		}
	}
	return true
}

func stringsEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

// Write licenses-small.gob, which is the same as licenses.gob except
// that license contents are limited to 256 bytes.
func writeSmallLicenseFile() error {
	f, err := os.Open("licenses.gob")
	if err != nil {
		return err
	}
	d := gob.NewDecoder(f)
	var ld LicenseData
	if err := d.Decode(&ld); err != nil {
		return err
	}
	for _, c := range ld.Contents {
		if len(c.Contents) > 256 {
			c.Contents = c.Contents[:256]
		}
	}
	out, err := os.Create("licenses-small.gob")
	if err != nil {
		return err
	}
	if err := gob.NewEncoder(out).Encode(&ld); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
