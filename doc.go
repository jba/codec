// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package codec implements an encoder for Go values. It relies on code generation
rather than reflection, so it is significantly faster than reflection-based
encoders like gob. It can also preserve sharing among pointers (but not other
forms of sharing, like sub-slices).


Struct Tags

Struct tags in the style of encoding/json are supported, under the name "codec":

    type T struct {
        A int `codec:"B"`
        C int `codec:"-"`
    }

Here, field A will use the name "B" and field C will be omitted. There is no need
for the omitempty option because the encoding always omits zero values.

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

XXXXXXXXXXXXXXXX

TODO:
   - how we read old code to get field names



*/
package codec
