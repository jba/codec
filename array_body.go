// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

const arrayBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
« $elTypeID := typeID .Type.Elem »
« $elTypeCodec := print $elTypeID "_codec" »
« $sliceTypeID := typeID .SliceType »
« $sliceTypeCodec := print $sliceTypeID "_codec" »

var «$typeID»_type = reflect.TypeOf((*«$goName»)(nil)).Elem()

type «$typeName» struct {
	codecapi.NonStruct
	«if .ElField -»
		«$elTypeCodec» *«$elTypeCodec»
	«end -»
	«if not .IsBytes -»
		«$sliceTypeCodec» *«$sliceTypeCodec»
	«end -»
}

func (c *«$typeName») TypesUsed() []reflect.Type {
	return []reflect.Type{
		«- if .ElField»
			«$elTypeID»_type,
		«end -»
		«if not .IsBytes» «$sliceTypeID»_type, «end»
	}
}

func (c *«$typeName») SetCodecs(tcs []codecapi.TypeCodec) {
	« $i := 0 -»
	«if .ElField -»
		c.«$elTypeCodec» = tcs[0].(*«$elTypeCodec»)
		« $i = 1 -»
	«end -»
	«if not .IsBytes -»
		c.«$sliceTypeID»_codec = tcs[«$i»].(*«$sliceTypeID»_codec)
	«end -»
}

func (c *«$typeName») Encode(e *codecapi.Encoder, x interface{}) {
	a := x.(«$goName»)
	c.encode(e, &a)
}

func (c *«$typeName») encode(e *codecapi.Encoder, s *«$goName») {
	«encodeStmt .SliceType "(*s)[:]"»
}

func (c *«$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c *«$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	«if .IsBytes -»
		b := d.DecodeBytes()
		copy((*p)[:], b)
	«else -»
		n := d.StartList()
		if n < 0 { return }
	    if n != «.Type.Len» {
			codecapi.Failf("array size mismatch: got %d, want «.Type.Len»", n)
		}
		for i := 0; i < n; i++ {
			«decodeStmt .Type.Elem "(*p)[i]"»
		}
	«end -»
}

func init() {
  codecapi.Register(«$goName»{}, func() codecapi.TypeCodec { return &«$typeName»{} })
}
`
