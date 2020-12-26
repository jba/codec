package main

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/google/licensecheck"
)

var (
	licenses      = licenseBenchmarkData("licenses")
	licensesSmall = licenseBenchmarkData("licenses-small")
)

func licenseBenchmarkData(name string) benchmarkData {
	return benchmarkData{
		name,
		func() (interface{}, error) {
			var ld LicenseData
			d, err := gobDecodeFile(name+".gob", &ld)
			if err == nil {
				n := 0
				n100 := 0
				for _, c := range ld.Contents {
					cvgFloats(c.OldCoverage, &n, &n100)
					cvgFloats(c.NewCoverage, &n, &n100)
				}
				fmt.Println(n, n100)
			}
			return d, err
		},
		func() interface{} { return new(*LicenseData) },
	}
}

func cvgFloats(c licensecheck.Coverage, n, n100 *int) {
	(*n) += 1 + len(c.Match)
	if c.Percent == 100 {
		(*n100)++
	}
	for _, m := range c.Match {
		if m.Percent == 100 {
			(*n100)++
		}
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

// LicenseContents hold the contents of a license file, and information derived from it.
type LicenseContents struct {
	Contents     []byte
	ContentsHash [sha256.Size]byte     // SHA256 of the contents, to dedup equal contents
	OldTypes     []string              // from the DB, stored in the gob file
	OldCoverage  licensecheck.Coverage // ditto
	NewTypes     []string              // not populated from the gob
	NewCoverage  licensecheck.Coverage // ditto
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
