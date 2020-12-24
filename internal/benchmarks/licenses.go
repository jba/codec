package main

import (
	"crypto/sha256"
	"encoding/gob"
	"os"

	"github.com/google/licensecheck"
)

var (
	licenses = benchmarkData{
		"licenses",
		func() (interface{}, error) { var ld LicenseData; return gobDecodeFile("licenses.gob", &ld) },
		func() interface{} { return new(*LicenseData) },
	}
	licensesSmall = benchmarkData{
		"licenses-small",
		func() (interface{}, error) { var ld LicenseData; return gobDecodeFile("licenses-small.gob", &ld) },
		func() interface{} { return new(*LicenseData) },
	}
)

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
