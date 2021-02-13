// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"reflect"

	"github.com/jba/codec/codecapi"
)

var array_1_int_type = reflect.TypeOf((*[1]int)(nil)).Elem()

type array_1_int_codec struct {
	slice_int_codec *slice_int_codec
}

func (c *array_1_int_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {

}

func (c *array_1_int_codec) Fields() []string { return nil }

func (c *array_1_int_codec) TypesUsed() []reflect.Type {
	return []reflect.Type{

		slice_int_type,
	}
}

func (c *array_1_int_codec) CodecsUsed(tcs []codecapi.TypeCodec) {
	c.slice_int_codec = tcs[0].(*slice_int_codec)
}

func (c *array_1_int_codec) Encode(e *codecapi.Encoder, x interface{}) {
	a := x.([1]int)
	c.encode(e, &a)
}

func (c *array_1_int_codec) encode(e *codecapi.Encoder, s *[1]int) {
	c.slice_int_codec.encode(e, (*s)[:])
}

func (c *array_1_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x [1]int
	c.decode(d, &x)
	return x
}

func (c *array_1_int_codec) decode(d *codecapi.Decoder, p *[1]int) {
	n := d.StartList()
	if n < 0 {
		return
	}
	if n != 1 {
		codecapi.Failf("array size mismatch: got %d, want 1", n)
	}
	for i := 0; i < n; i++ {
		(*p)[i] = int(d.DecodeInt())
	}
}

func init() {
	codecapi.Register([1]int{}, func() codecapi.TypeCodec { return &array_1_int_codec{} })
}

var slice_int_type = reflect.TypeOf((*[]int)(nil)).Elem()

type slice_int_codec struct {
}

func (c *slice_int_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {

}

func (c *slice_int_codec) Fields() []string { return nil }

func (c *slice_int_codec) TypesUsed() []reflect.Type {
	return nil
}

func (c *slice_int_codec) CodecsUsed(tcs []codecapi.TypeCodec) {
}

func (c *slice_int_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.([]int)) }

func (c *slice_int_codec) encode(e *codecapi.Encoder, s []int) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c *slice_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x []int
	c.decode(d, &x)
	return x
}

func (c *slice_int_codec) decode(d *codecapi.Decoder, p *[]int) {
	n := d.StartList()
	if n < 0 {
		return
	}
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = int(d.DecodeInt())
	}
	*p = s
}

func init() {
	codecapi.Register([]int(nil), func() codecapi.TypeCodec { return &slice_int_codec{} })
}

var ptr_smallStruct_type = reflect.TypeOf((*smallStruct)(nil))

type ptr_smallStruct_codec struct {
	smallStruct_codec *smallStruct_codec
}

func (c *ptr_smallStruct_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {
	c.smallStruct_codec = tcs[reflect.TypeOf((*smallStruct)(nil)).Elem()].(*smallStruct_codec)
}

func (c ptr_smallStruct_codec) Fields() []string { return nil }

func (c ptr_smallStruct_codec) TypesUsed() []reflect.Type { return []reflect.Type{smallStruct_type} }

func (c *ptr_smallStruct_codec) CodecsUsed(tcs []codecapi.TypeCodec) {
	c.smallStruct_codec = tcs[0].(*smallStruct_codec)
}

func (c ptr_smallStruct_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(*smallStruct))
}

func (c ptr_smallStruct_codec) encode(e *codecapi.Encoder, x *smallStruct) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	c.smallStruct_codec.encode(e, x)
}

func (c ptr_smallStruct_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *smallStruct
	c.decode(d, &x)
	return x
}

func (c ptr_smallStruct_codec) decode(d *codecapi.Decoder, p **smallStruct) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*smallStruct)
		return
	}
	var x smallStruct
	d.StoreRef(&x)
	c.smallStruct_codec.decode(d, &x)
	*p = &x
}

var smallStruct_type = ptr_smallStruct_type.Elem()

type smallStruct_codec struct {
	fieldMap []int
}

func (c *smallStruct_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, fieldMap []int) {
	c.fieldMap = fieldMap
}

func (c *smallStruct_codec) Fields() []string {
	return []string{"X"}
}

func (c *smallStruct_codec) TypesUsed() []reflect.Type {
	return []reflect.Type{}
}

func (c *smallStruct_codec) CodecsUsed(tcs []codecapi.TypeCodec) {
}

func (c *smallStruct_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(smallStruct)
	c.encode(e, &s)
}

func (c *smallStruct_codec) encode(e *codecapi.Encoder, x *smallStruct) {
	e.StartStruct()
	if x.X != 0 {
		e.EncodeUint(0)
		e.EncodeInt(int64(x.X))
	}
	e.EndStruct()
}

func (c *smallStruct_codec) Decode(d *codecapi.Decoder) interface{} {
	var x smallStruct
	c.decode(d, &x)
	return x
}

func (c *smallStruct_codec) decode(d *codecapi.Decoder, x *smallStruct) {
	d.StartStruct()
	for {
		n := d.NextStructField(c.fieldMap)
		if n == -1 {
			break
		}
		switch n {
		case 0:
			x.X = int(d.DecodeInt())
		default:
			d.UnknownField("smallStruct", n)
		}
	}
}

func init() {
	codecapi.Register(smallStruct{}, func() codecapi.TypeCodec { return &smallStruct_codec{} })
	codecapi.Register(&smallStruct{}, func() codecapi.TypeCodec { return &ptr_smallStruct_codec{} })
}

var map__1_int_smallStruct_type = reflect.TypeOf((*map[[1]int]smallStruct)(nil)).Elem()

type map__1_int_smallStruct_codec struct {
	array_1_int_codec *array_1_int_codec
	smallStruct_codec *smallStruct_codec
}

func (c *map__1_int_smallStruct_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {
	c.array_1_int_codec = tcs[reflect.TypeOf((*[1]int)(nil)).Elem()].(*array_1_int_codec)
	c.smallStruct_codec = tcs[reflect.TypeOf((*smallStruct)(nil)).Elem()].(*smallStruct_codec)
}

func (c *map__1_int_smallStruct_codec) Fields() []string { return nil }

func (c *map__1_int_smallStruct_codec) TypesUsed() []reflect.Type {
	// TODO:  generate a slice literal
	var types []reflect.Type
	types = append(types, array_1_int_type)
	types = append(types, smallStruct_type)
	return types
}

func (c *map__1_int_smallStruct_codec) CodecsUsed(tcs []codecapi.TypeCodec) {
	c.array_1_int_codec = tcs[0].(*array_1_int_codec)
	c.smallStruct_codec = tcs[1].(*smallStruct_codec)
}

func (c *map__1_int_smallStruct_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(map[[1]int]smallStruct))
}

func (c *map__1_int_smallStruct_codec) encode(e *codecapi.Encoder, m map[[1]int]smallStruct) {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2 * len(m))
	for k, v := range m {
		c.array_1_int_codec.encode(e, &k)
		c.smallStruct_codec.encode(e, &v)
	}
}

func (c *map__1_int_smallStruct_codec) Decode(d *codecapi.Decoder) interface{} {
	var x map[[1]int]smallStruct
	c.decode(d, &x)
	return x
}

func (c *map__1_int_smallStruct_codec) decode(d *codecapi.Decoder, p *map[[1]int]smallStruct) {
	n2 := d.StartList()
	if n2 < 0 {
		return
	}
	n := n2 / 2
	m := make(map[[1]int]smallStruct, n)
	var k [1]int
	var v smallStruct
	for i := 0; i < n; i++ {
		c.array_1_int_codec.decode(d, &k)
		c.smallStruct_codec.decode(d, &v)
		m[k] = v
	}
	*p = m
}

func init() {
	codecapi.Register(map[[1]int]smallStruct(nil), func() codecapi.TypeCodec { return &map__1_int_smallStruct_codec{} })
}
