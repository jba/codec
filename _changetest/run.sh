#!/bin/sh -e

rm -f change.gen.go

go run -tags version1,gen .
go run -tags version1,encode .
go run -tags version2,gen .
go run -tags version2,decode .

rm change.gen.go
rm change.enc
