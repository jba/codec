// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package codec implements an encoder for Go values. It relies on code generation
rather than reflection, so it is significantly faster than reflection-based
encoders like gob. It can also preserve sharing among pointers (but not other
forms of sharing, like sub-slices).


Generating Code

The package supports built-in types (int, string and so on) out of the box,
but for any other type you must generate code by calling GenerateFile.
This can be done with a small program in your project's directory:

    // file generate.go
    //+build ignore

	package main

	import (
	   "myproject"
	   "github.com/jba/codec"
	)

	func main() {
		err := codec.GenerateFile("types.gen.go", "mypkg", nil, []Type1{}, Type2{})
		if err != nil {
			log.Fatal(err)
		}
	}

The "//+build ignore" tag prevents the program from being compiled as part of your package.
Instead, invoke it directly with "go run". Use "go generate" to do so if you like:

    //go:generate go run generate.go

On subsequent runs, the generator reads the generated file to get the names and
order of all struct fields. It uses this information to generate correct code
when fields are moved or added. Make sure the old generated files remain available
to the generator, or changes to your structs may result in existing encoded data
being decoded incorrectly.


Encoding and Decoding

Create an Encoder by passing it an io.Writer:

	var buf bytes.Buffer
	e := codec.NewEncoder(&buf, nil)

Then use it to encode one or more values:

   if err := e.Encode(x); err != nil { ... }

To decode, pass an io.Reader to NewDecoder, and call Decode:

   f, err := os.Open(filename)
   ...
   d := codec.NewDecoder(f, nil)
   value, err := d.Decode()
   ...


Struct Tags

Struct tags in the style of encoding/json are supported, under the name "codec":

    type T struct {
        A int `codec:"B"`
        C int `codec:"-"`
    }

Here, field A will use the name "B" and field C will be omitted. There is no need
for the omitempty option because the encoder always omits zero values.

Since the encoding uses numbers for fields instead of names, renaming a field doesn't
actually affect the encoding. It does matter if subsequent changes are made to the struct,
however. For example, say that originally T was

    type T struct {
        B int
    }

but now you'd like to rename the field to "A":

    type T struct {
        A int
    }

The generator will treat "A" as a new field. Use a tag to express that it is a renaming:

    type T struct {
        A int `codec:"B"`
    }

*/
package codec
