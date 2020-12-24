// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	
	"github.com/jba/codec/codecapi"
	"go/token"
	"time"
)

// Fields of generatedTestTypes: Node Slice Array Map Struct Time DefSlice DefArray DefMap Pos Stocks

type ptr_generatedTestTypes_codec struct{}

func (ptr_generatedTestTypes_codec) Init() {}

func (c ptr_generatedTestTypes_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(*generatedTestTypes))
}

func (c ptr_generatedTestTypes_codec) encode(e *codecapi.Encoder, x *generatedTestTypes) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(generatedTestTypes_codec{}).encode(e, x)
}

func (c ptr_generatedTestTypes_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *generatedTestTypes
	c.decode(d, &x)
	return x
}

func (c ptr_generatedTestTypes_codec) decode(d *codecapi.Decoder, p **generatedTestTypes) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*generatedTestTypes)
		return
	}
	var x generatedTestTypes
	d.StoreRef(&x)
	(generatedTestTypes_codec{}).decode(d, &x)
	*p = &x
}

type generatedTestTypes_codec struct{}

func (generatedTestTypes_codec) Init() {}

func (c generatedTestTypes_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(generatedTestTypes)
	c.encode(e, &s)
}

func (c generatedTestTypes_codec) encode(e *codecapi.Encoder, x *generatedTestTypes) {
	e.StartStruct()
	if x.Node != nil {
		e.EncodeUint(0)
		(ptr_node_codec{}).encode(e, x.Node)
	}
	if x.Slice != nil {
		e.EncodeUint(1)
		(slice_int_codec{}).encode(e, x.Slice)
	}

	e.EncodeUint(2)
	(array_1_int_codec{}).encode(e, &x.Array)
	if x.Map != nil {
		e.EncodeUint(3)
		(map_string_bool_codec{}).encode(e, x.Map)
	}

	e.EncodeUint(4)
	(structType_codec{}).encode(e, &x.Struct)

	e.EncodeUint(5)
	(time_Time_codec{}).encode(e, x.Time)
	if x.DefSlice != nil {
		e.EncodeUint(6)
		(definedSlice_codec{}).encode(e, definedSlice(x.DefSlice))
	}

	e.EncodeUint(7)
	(definedArray_codec{}).encode(e, &x.DefArray)
	if x.DefMap != nil {
		e.EncodeUint(8)
		(definedMap_codec{}).encode(e, definedMap(x.DefMap))
	}
	if x.Pos != 0 {
		e.EncodeUint(9)
		e.EncodeInt(int64(x.Pos))
	}
	if x.Stocks != nil {
		e.EncodeUint(10)
		(slice_StockData_codec{}).encode(e, x.Stocks)
	}
	e.EndStruct()
}

func (c generatedTestTypes_codec) Decode(d *codecapi.Decoder) interface{} {
	var x generatedTestTypes
	c.decode(d, &x)
	return x
}

func (c generatedTestTypes_codec) decode(d *codecapi.Decoder, x *generatedTestTypes) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			(ptr_node_codec{}).decode(d, &x.Node)
		case 1:
			(slice_int_codec{}).decode(d, &x.Slice)
		case 2:
			(array_1_int_codec{}).decode(d, &x.Array)
		case 3:
			(map_string_bool_codec{}).decode(d, &x.Map)
		case 4:
			(structType_codec{}).decode(d, &x.Struct)
		case 5:
			(time_Time_codec{}).decode(d, &x.Time)
		case 6:
			(definedSlice_codec{}).decode(d, &x.DefSlice)
		case 7:
			(definedArray_codec{}).decode(d, &x.DefArray)
		case 8:
			(definedMap_codec{}).decode(d, &x.DefMap)
		case 9:
			x.Pos = token.Pos(d.DecodeInt())
		case 10:
			(slice_StockData_codec{}).decode(d, &x.Stocks)
		default:
			d.UnknownField("generatedTestTypes", n)
		}
	}
}

func init() {
	codecapi.Register(generatedTestTypes{}, generatedTestTypes_codec{})
	codecapi.Register(&generatedTestTypes{}, ptr_generatedTestTypes_codec{})
}

// Fields of node: Value Next

type ptr_node_codec struct{}

func (ptr_node_codec) Init() {}

func (c ptr_node_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(*node)) }

func (c ptr_node_codec) encode(e *codecapi.Encoder, x *node) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(node_codec{}).encode(e, x)
}

func (c ptr_node_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *node
	c.decode(d, &x)
	return x
}

func (c ptr_node_codec) decode(d *codecapi.Decoder, p **node) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*node)
		return
	}
	var x node
	d.StoreRef(&x)
	(node_codec{}).decode(d, &x)
	*p = &x
}

type node_codec struct{}

func (node_codec) Init() {}

func (c node_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(node)
	c.encode(e, &s)
}

func (c node_codec) encode(e *codecapi.Encoder, x *node) {
	e.StartStruct()
	if x.Value != 0 {
		e.EncodeUint(0)
		e.EncodeInt(int64(x.Value))
	}
	if x.Next != nil {
		e.EncodeUint(1)
		(ptr_node_codec{}).encode(e, x.Next)
	}
	e.EndStruct()
}

func (c node_codec) Decode(d *codecapi.Decoder) interface{} {
	var x node
	c.decode(d, &x)
	return x
}

func (c node_codec) decode(d *codecapi.Decoder, x *node) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			x.Value = int(d.DecodeInt())
		case 1:
			(ptr_node_codec{}).decode(d, &x.Next)
		default:
			d.UnknownField("node", n)
		}
	}
}

func init() {
	codecapi.Register(node{}, node_codec{})
	codecapi.Register(&node{}, ptr_node_codec{})
}

type slice_int_codec struct{}

func (slice_int_codec) Init() {}

func (c slice_int_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.([]int)) }

func (c slice_int_codec) encode(e *codecapi.Encoder, s []int) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c slice_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x []int
	c.decode(d, &x)
	return x
}

func (c slice_int_codec) decode(d *codecapi.Decoder, p *[]int) {
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
	codecapi.Register([]int(nil), slice_int_codec{})
}

type array_1_int_codec struct{}

func (array_1_int_codec) Init() {}

func (c array_1_int_codec) Encode(e *codecapi.Encoder, x interface{}) {
	a := x.([1]int)
	c.encode(e, &a)
}

func (c array_1_int_codec) encode(e *codecapi.Encoder, s *[1]int) {
	(slice_int_codec{}).encode(e, (*s)[:])
}

func (c array_1_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x [1]int
	c.decode(d, &x)
	return x
}

func (c array_1_int_codec) decode(d *codecapi.Decoder, p *[1]int) {
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
	codecapi.Register([1]int{}, array_1_int_codec{})
}

type map_string_bool_codec struct{}

func (c map_string_bool_codec) Init() {}

func (c map_string_bool_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(map[string]bool))
}

func (c map_string_bool_codec) encode(e *codecapi.Encoder, m map[string]bool) {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2 * len(m))
	for k, v := range m {
		e.EncodeString(k)
		e.EncodeBool(v)
	}
}

func (c map_string_bool_codec) Decode(d *codecapi.Decoder) interface{} {
	var x map[string]bool
	c.decode(d, &x)
	return x
}

func (c map_string_bool_codec) decode(d *codecapi.Decoder, p *map[string]bool) {
	n2 := d.StartList()
	if n2 < 0 {
		return
	}
	n := n2 / 2
	m := make(map[string]bool, n)
	var k string
	var v bool
	for i := 0; i < n; i++ {
		k = d.DecodeString()
		v = d.DecodeBool()
		m[k] = v
	}
	*p = m
}

func init() { codecapi.Register(map[string]bool(nil), map_string_bool_codec{}) }

// Fields of structType: N

type ptr_structType_codec struct{}

func (ptr_structType_codec) Init() {}

func (c ptr_structType_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(*structType))
}

func (c ptr_structType_codec) encode(e *codecapi.Encoder, x *structType) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(structType_codec{}).encode(e, x)
}

func (c ptr_structType_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *structType
	c.decode(d, &x)
	return x
}

func (c ptr_structType_codec) decode(d *codecapi.Decoder, p **structType) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*structType)
		return
	}
	var x structType
	d.StoreRef(&x)
	(structType_codec{}).decode(d, &x)
	*p = &x
}

type structType_codec struct{}

func (structType_codec) Init() {}

func (c structType_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(structType)
	c.encode(e, &s)
}

func (c structType_codec) encode(e *codecapi.Encoder, x *structType) {
	e.StartStruct()

	e.EncodeUint(0)
	(node_codec{}).encode(e, &x.N)
	e.EndStruct()
}

func (c structType_codec) Decode(d *codecapi.Decoder) interface{} {
	var x structType
	c.decode(d, &x)
	return x
}

func (c structType_codec) decode(d *codecapi.Decoder, x *structType) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			(node_codec{}).decode(d, &x.N)
		default:
			d.UnknownField("structType", n)
		}
	}
}

func init() {
	codecapi.Register(structType{}, structType_codec{})
	codecapi.Register(&structType{}, ptr_structType_codec{})
}

type time_Time_codec struct{}

func (c time_Time_codec) Init() {}

func (c time_Time_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(time.Time)) }

func (c time_Time_codec) encode(e *codecapi.Encoder, m time.Time) {
	data, err := m.MarshalBinary()
	if err != nil {
		codecapi.Fail(err)
	}
	e.EncodeBytes(data)
}

func (c time_Time_codec) Decode(d *codecapi.Decoder) interface{} {
	var x time.Time
	c.decode(d, &x)
	return x
}

func (c time_Time_codec) decode(d *codecapi.Decoder, p *time.Time) {
	data := d.DecodeBytes()
	if err := p.UnmarshalBinary(data); err != nil {
		codecapi.Fail(err)
	}
}

func init() { codecapi.Register(*new(time.Time), time_Time_codec{}) }

type definedSlice_codec struct{}

func (definedSlice_codec) Init() {}

func (c definedSlice_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(definedSlice)) }

func (c definedSlice_codec) encode(e *codecapi.Encoder, s definedSlice) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c definedSlice_codec) Decode(d *codecapi.Decoder) interface{} {
	var x definedSlice
	c.decode(d, &x)
	return x
}

func (c definedSlice_codec) decode(d *codecapi.Decoder, p *definedSlice) {
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
	codecapi.Register(definedSlice(nil), definedSlice_codec{})
}

type definedArray_codec struct{}

func (definedArray_codec) Init() {}

func (c definedArray_codec) Encode(e *codecapi.Encoder, x interface{}) {
	a := x.(definedArray)
	c.encode(e, &a)
}

func (c definedArray_codec) encode(e *codecapi.Encoder, s *definedArray) {
	(slice_int_codec{}).encode(e, (*s)[:])
}

func (c definedArray_codec) Decode(d *codecapi.Decoder) interface{} {
	var x definedArray
	c.decode(d, &x)
	return x
}

func (c definedArray_codec) decode(d *codecapi.Decoder, p *definedArray) {
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
	codecapi.Register(definedArray{}, definedArray_codec{})
}

type definedMap_codec struct{}

func (c definedMap_codec) Init() {}

func (c definedMap_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(definedMap)) }

func (c definedMap_codec) encode(e *codecapi.Encoder, m definedMap) {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2 * len(m))
	for k, v := range m {
		e.EncodeString(k)
		e.EncodeBool(v)
	}
}

func (c definedMap_codec) Decode(d *codecapi.Decoder) interface{} {
	var x definedMap
	c.decode(d, &x)
	return x
}

func (c definedMap_codec) decode(d *codecapi.Decoder, p *definedMap) {
	n2 := d.StartList()
	if n2 < 0 {
		return
	}
	n := n2 / 2
	m := make(definedMap, n)
	var k string
	var v bool
	for i := 0; i < n; i++ {
		k = d.DecodeString()
		v = d.DecodeBool()
		m[k] = v
	}
	*p = m
}

func init() { codecapi.Register(definedMap(nil), definedMap_codec{}) }

type slice_StockData_codec struct{}

func (slice_StockData_codec) Init() {}

func (c slice_StockData_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.([]StockData))
}

func (c slice_StockData_codec) encode(e *codecapi.Encoder, s []StockData) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		(StockData_codec{}).encode(e, &x)
	}
}

func (c slice_StockData_codec) Decode(d *codecapi.Decoder) interface{} {
	var x []StockData
	c.decode(d, &x)
	return x
}

func (c slice_StockData_codec) decode(d *codecapi.Decoder, p *[]StockData) {
	n := d.StartList()
	if n < 0 {
		return
	}
	s := make([]StockData, n)
	for i := 0; i < n; i++ {
		(StockData_codec{}).decode(d, &s[i])
	}
	*p = s
}

func init() {
	codecapi.Register([]StockData(nil), slice_StockData_codec{})
}

// Fields of StockData: Symbol Intervals

type ptr_StockData_codec struct{}

func (ptr_StockData_codec) Init() {}

func (c ptr_StockData_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(*StockData)) }

func (c ptr_StockData_codec) encode(e *codecapi.Encoder, x *StockData) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(StockData_codec{}).encode(e, x)
}

func (c ptr_StockData_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *StockData
	c.decode(d, &x)
	return x
}

func (c ptr_StockData_codec) decode(d *codecapi.Decoder, p **StockData) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*StockData)
		return
	}
	var x StockData
	d.StoreRef(&x)
	(StockData_codec{}).decode(d, &x)
	*p = &x
}

type StockData_codec struct{}

func (StockData_codec) Init() {}

func (c StockData_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(StockData)
	c.encode(e, &s)
}

func (c StockData_codec) encode(e *codecapi.Encoder, x *StockData) {
	e.StartStruct()
	if x.Symbol != "" {
		e.EncodeUint(0)
		e.EncodeString(x.Symbol)
	}
	if x.Intervals != nil {
		e.EncodeUint(1)
		(slice_Interval_codec{}).encode(e, x.Intervals)
	}
	e.EndStruct()
}

func (c StockData_codec) Decode(d *codecapi.Decoder) interface{} {
	var x StockData
	c.decode(d, &x)
	return x
}

func (c StockData_codec) decode(d *codecapi.Decoder, x *StockData) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			x.Symbol = d.DecodeString()
		case 1:
			(slice_Interval_codec{}).decode(d, &x.Intervals)
		default:
			d.UnknownField("StockData", n)
		}
	}
}

func init() {
	codecapi.Register(StockData{}, StockData_codec{})
	codecapi.Register(&StockData{}, ptr_StockData_codec{})
}

type slice_Interval_codec struct{}

func (slice_Interval_codec) Init() {}

func (c slice_Interval_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.([]Interval)) }

func (c slice_Interval_codec) encode(e *codecapi.Encoder, s []Interval) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		(Interval_codec{}).encode(e, &x)
	}
}

func (c slice_Interval_codec) Decode(d *codecapi.Decoder) interface{} {
	var x []Interval
	c.decode(d, &x)
	return x
}

func (c slice_Interval_codec) decode(d *codecapi.Decoder, p *[]Interval) {
	n := d.StartList()
	if n < 0 {
		return
	}
	s := make([]Interval, n)
	for i := 0; i < n; i++ {
		(Interval_codec{}).decode(d, &s[i])
	}
	*p = s
}

func init() {
	codecapi.Register([]Interval(nil), slice_Interval_codec{})
}

// Fields of Interval: Start End Open Close Low High Volume

type ptr_Interval_codec struct{}

func (ptr_Interval_codec) Init() {}

func (c ptr_Interval_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(*Interval)) }

func (c ptr_Interval_codec) encode(e *codecapi.Encoder, x *Interval) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(Interval_codec{}).encode(e, x)
}

func (c ptr_Interval_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *Interval
	c.decode(d, &x)
	return x
}

func (c ptr_Interval_codec) decode(d *codecapi.Decoder, p **Interval) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*Interval)
		return
	}
	var x Interval
	d.StoreRef(&x)
	(Interval_codec{}).decode(d, &x)
	*p = &x
}

type Interval_codec struct{}

func (Interval_codec) Init() {}

func (c Interval_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(Interval)
	c.encode(e, &s)
}

func (c Interval_codec) encode(e *codecapi.Encoder, x *Interval) {
	e.StartStruct()

	e.EncodeUint(0)
	(time_Time_codec{}).encode(e, x.Start)

	e.EncodeUint(1)
	(time_Time_codec{}).encode(e, x.End)
	if x.Open != 0 {
		e.EncodeUint(2)
		e.EncodeFloat(x.Open)
	}
	if x.Close != 0 {
		e.EncodeUint(3)
		e.EncodeFloat(x.Close)
	}
	if x.Low != 0 {
		e.EncodeUint(4)
		e.EncodeFloat(x.Low)
	}
	if x.High != 0 {
		e.EncodeUint(5)
		e.EncodeFloat(x.High)
	}
	if x.Volume != 0 {
		e.EncodeUint(6)
		e.EncodeFloat(x.Volume)
	}
	e.EndStruct()
}

func (c Interval_codec) Decode(d *codecapi.Decoder) interface{} {
	var x Interval
	c.decode(d, &x)
	return x
}

func (c Interval_codec) decode(d *codecapi.Decoder, x *Interval) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			(time_Time_codec{}).decode(d, &x.Start)
		case 1:
			(time_Time_codec{}).decode(d, &x.End)
		case 2:
			x.Open = d.DecodeFloat()
		case 3:
			x.Close = d.DecodeFloat()
		case 4:
			x.Low = d.DecodeFloat()
		case 5:
			x.High = d.DecodeFloat()
		case 6:
			x.Volume = d.DecodeFloat()
		default:
			d.UnknownField("Interval", n)
		}
	}
}

func init() {
	codecapi.Register(Interval{}, Interval_codec{})
	codecapi.Register(&Interval{}, ptr_Interval_codec{})
}
