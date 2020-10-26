#!/bin/bash

filename=types.gen_test.go

rm -f $filename
go test -generate $filename
sed -i \
  -e 's/codec\.//g' \
  -e 's@"github.com/jba/codec"@@' \
  $filename
