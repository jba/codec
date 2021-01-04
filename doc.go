// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package codec implements an encoder for Go values. It relies on code generation
rather than reflection, so it is significantly faster than reflection-based
encoders like gob. It can also preserve sharing among pointers (but not other
forms of sharing, like sub-slices).

Encodings with maps are not deterministic, due to the non-deterministic order of
map iteration.


Generating Code

The package supports Go built-in types (int, string and so on) out of the box,
but for any other type you must generate code by calling GenerateFile. This can
be done with a small program in your project's directory:

    // file generate.go
    //+build ignore

	package main

	import (
	   "mypkg"
	   "github.com/jba/codec"
	)

	func main() {
		err := codec.GenerateFile("types.gen.go", "mypkg", nil,
			[]mypkg.Type1{}, &mypkg.Type2{})
		if err != nil {
			log.Fatal(err)
		}
	}

Code will be generated for each type listed and for all types they contain. So
this program will generate code for []mypkg.Type1, mypkg.Type1, *mypkg.Type2,
and mypkg.Type2.

The "//+build ignore" tag prevents the program from being compiled as part of
your package. Instead, invoke it directly with "go run". Use "go generate" to do
so if you like:

    //go:generate go run generate.go

On subsequent runs, the generator reads the generated file to get the names and
order of all struct fields. It uses this information to generate correct code
when fields are moved or added. Make sure the old generated files remain
available to the generator, or changes to your structs may result in existing
encoded data being decoded incorrectly.


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


Sharing and Cycles

By default, if two pointers point to the same value, that value will be
duplicated upon decoding. If there is a cycle, where a value directly or
indirectly points to itself, then the encoder will crash by exceeding available
stack space. This is the same behavior as encoding/gob and many other encoders.

Set EncodeOptions.TrackPointers to true to preserve pointer sharing and cycles,
at the cost of slower encoding.

Other forms of memory sharing are not preserved. For example, if two slices
refer to the same underlying array during encoding, they will refer to separate
arrays after decoding.


Struct Tags

Struct tags in the style of encoding/json are supported, under the name "codec".
You can easily generate code for structs designed for the encoding/json package
by changing the name to "json" in an option to GenerateFile.

An example:

    type T struct {
        A int `codec:"B"`
        C int `codec:"-"`
    }

Here, field A will use the name "B" and field C will be omitted. There is no
need for the omitempty option because the encoder always omits zero values.

Since the encoding uses numbers for fields instead of names, renaming a field
doesn't actually affect the encoding. It does matter if subsequent changes are
made to the struct, however. For example, say that originally T was

    type T struct {
        A int
    }

but you rename the field to "B":

    type T struct {
        B int
    }

The generator will treat "B" as a new field. Data encoded with "A" will not be
decoded into "B". So you should use a tag to express that it is a renaming:

    type T struct {
        B int `codec:"A"`
    }

*/
package codec
