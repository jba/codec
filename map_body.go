// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

const mapBody = `
«/*»
Template body for a map type.
A nil map is encoded as a zero.
A map of size N is encoded as a list of length 2N, containing alternating
keys and values.

In the decode function, we declare a variable v to hold the decoded map value
rather than decoding directly into m[v]. This is necessary for decode
functions that take pointers: you can't take a pointer to a map element.
«*/»

« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
« $keyTypeID := typeID .Type.Key »
« $elTypeID := typeID .Type.Elem »

var «$typeID»_type = reflect.TypeOf((*«$goName»)(nil)).Elem()

type «$typeName» struct {
	codecapi.NonStruct
	«if .KeyField -»
		«$keyTypeID»_codec *«$keyTypeID»_codec
	«end -»
	«if .ElField -»
		«$elTypeID»_codec *«$elTypeID»_codec
	«end -»
}

func (c *«$typeName») TypesUsed() []reflect.Type {
	return []reflect.Type{
		«if .KeyField» «$keyTypeID»_type, «end»
		«if .ElField»  «$elTypeID»_type,  «end»
	}
}

func (c *«$typeName») SetCodecs(tcs []codecapi.TypeCodec) {
	«if .KeyField -»
		c.«$keyTypeID»_codec = tcs[0].(*«$keyTypeID»_codec)
	«end -»
	«if .ElField -»
		c.«$elTypeID»_codec = tcs[«if .KeyField»1«else»0«end»].(*«$elTypeID»_codec)
	«end -»
}

func (c *«$typeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(«$goName»)) }

func (c *«$typeName») encode(e *codecapi.Encoder, m «$goName») {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2*len(m))
	for k, v := range m {
		«encodeStmt .Type.Key "k"»
		«encodeStmt .Type.Elem "v"»
	}
}

func (c *«$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c *«$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	n2 := d.StartList()
	if n2 < 0 { return }
	n := n2/2
	m := make(«$goName», n)
	var k «goName .Type.Key»
	var v «goName .Type.Elem»
	for i := 0; i < n; i++ {
		«decodeStmt .Type.Key "k"»
		«decodeStmt .Type.Elem "v"»
		m[k] = v
	}
	*p = m
}

func init() { codecapi.Register(«$goName»(nil), func() codecapi.TypeCodec { return &«$typeName»{} }) }
`