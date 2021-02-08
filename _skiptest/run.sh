#!/bin/sh -e

go run -tags version1,gen .
go run -tags version1,encode .

rm skip.gen.go

go run -tags version2,gen .
go run -tags version2,decode .

rm skip.gen.go
rm skip.enc
