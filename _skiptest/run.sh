#!/bin/sh -ex

go run -tags version1 gen.go skip1.go

go run -tags version1 encode.go skip1.go skip.gen.go

go run -tags version2 gen.go skip2.go

go run -tags version2 decode.go skip2.go skip.gen.go

rm skip.gen.go
rm skip.enc
