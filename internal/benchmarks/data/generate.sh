#!/bin/sh

rm -f *.gen.go
cat licenses-gob-* > licenses.gob
go run generate.go code ast stocks scores licenses
codecgen -o hyperledger.ugorji.gen.go hyperledger.go
codecgen -o licenses.ugorji.gen.go licenses.go
codecgen -o scores.ugorji.gen.go scores.go
codecgen -o stocks.ugorji.gen.go stocks.go
