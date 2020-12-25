// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"go/format"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/jba/codec/codecapi"
)

// GenerateFile writes encoders and decoders to filename.
// It generates code for the type of each given value, as well
// as any types they depend on.
// packageName is the name following the file's package declaration.
func GenerateFile(filename, packageName string, values ...interface{}) error {
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}
	fieldNames, err := readFieldNames(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := generate(f, packageName, fieldNames, values...); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// readFieldNames scan filename, if it exists, to get the previous field names for structs.
// It returns a map from struct name to list of field names.
func readFieldNames(filename string) (map[string][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	m := map[string][]string{}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "// Fields of ") {
			// form of line: // Fields of STRUCTNAME: NAME1 NAME2 ...
			parts := strings.Fields(line)
			structName := parts[3][:len(parts[3])-1] // remove final colon
			m[structName] = parts[4:]
		}
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func generate(w io.Writer, packageName string, fieldNames map[string][]string, vs ...interface{}) error {
	g := &generator{
		pkg:        packageName,
		done:       map[reflect.Type]bool{},
		fieldNames: fieldNames,
	}
	if g.fieldNames == nil {
		g.fieldNames = map[string][]string{}
	}
	funcs := template.FuncMap{
		"typeID":     g.typeIdentifier,
		"goName":     g.goName,
		"encodeStmt": g.encodeStatement,
		"decodeStmt": g.decodeStatement,
		"encodeFunc": g.encodeFunc,
		"decodeFunc": g.decodeFunc,
	}

	newTemplate := func(name, body string) *template.Template {
		return template.Must(template.New(name).Delims("«", "»").Funcs(funcs).Parse(body))
	}

	g.initialTemplate = newTemplate("initial", initialBody)
	g.sliceTemplate = newTemplate("slice", sliceBody)
	g.arrayTemplate = newTemplate("array", arrayBody)
	g.mapTemplate = newTemplate("map", mapBody)
	g.structTemplate = newTemplate("struct", structBody)
	g.marshalTemplate = newTemplate("marshaller", marshalBody)

	for _, v := range vs {
		g.todo = append(g.todo, reflect.TypeOf(v))
	}

	// Mark the built-in types as done.
	for _, t := range codecapi.BuiltinTypes {
		g.done[t] = true
		if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			g.done[t.Elem()] = true
		}
	}
	// The empty interface doesn't need any additional code. It's tricky to get
	// its reflect.Type: we need to dereference the pointer type.
	var iface interface{}
	g.done[reflect.TypeOf(&iface).Elem()] = true

	src, err := g.generate()
	if err != nil {
		return err
	}
	src, err = format.Source(src)
	if err != nil {
		return fmt.Errorf("format.Source: %v", err)
	}
	_, err = w.Write(src)
	return err
}

type generator struct {
	pkg             string
	todo            []reflect.Type
	done            map[reflect.Type]bool
	fieldNames      map[string][]string
	importMap       map[string]bool
	initialTemplate *template.Template
	sliceTemplate   *template.Template
	arrayTemplate   *template.Template
	mapTemplate     *template.Template
	structTemplate  *template.Template
	marshalTemplate *template.Template
}

func (g *generator) generate() ([]byte, error) {
	g.importMap = map[string]bool{"github.com/jba/codec/codecapi": true}
	var code []byte
	for len(g.todo) > 0 {
		t := g.todo[0]
		g.todo = g.todo[1:]
		if !g.done[t] {
			if t.PkgPath() != "" && t.PkgPath() != "main" {
				g.importMap[t.PkgPath()] = true
			}
			piece, err := g.gen(t)
			if err != nil {
				return nil, err
			}
			if piece != nil {
				code = append(code, piece...)
			}
			// We use the same code for T and *T, so both are done.
			g.done[t] = true
			g.done[reflect.PtrTo(t)] = true
		}
	}

	var imports []string
	for i := range g.importMap {
		imports = append(imports, i)
	}
	sort.Strings(imports)
	initial, err := execute(g.initialTemplate, struct {
		Package string
		Imports []string
	}{
		Package: g.pkg,
		Imports: imports,
	})
	if err != nil {
		return nil, err
	}
	return append(initial, code...), nil
}

var (
	binaryMarshalerType   = reflect.TypeOf(new(encoding.BinaryMarshaler)).Elem()
	binaryUnmarshalerType = reflect.TypeOf(new(encoding.BinaryUnmarshaler)).Elem()
	textMarshalerType     = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
	textUnmarshalerType   = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

func (g *generator) gen(t reflect.Type) ([]byte, error) {
	if t.Implements(binaryMarshalerType) && reflect.PtrTo(t).Implements(binaryUnmarshalerType) {
		return g.genMarshaler(t, "Binary")
	}
	if t.Implements(textMarshalerType) && reflect.PtrTo(t).Implements(textUnmarshalerType) {
		return g.genMarshaler(t, "Text")
	}

	switch t.Kind() {
	case reflect.Slice:
		return g.genSlice(t)
	case reflect.Array:
		return g.genArray(t)
	case reflect.Map:
		return g.genMap(t)
	case reflect.Struct:
		return g.genStruct(t)
	case reflect.Ptr:
		return g.gen(t.Elem())
	}
	return nil, nil
}

func (g *generator) genSlice(t reflect.Type) ([]byte, error) {
	et := t.Elem()
	g.todo = append(g.todo, et)
	return execute(g.sliceTemplate, struct {
		Type, ElType reflect.Type
	}{
		Type:   t,
		ElType: et,
	})
}

func (g *generator) genArray(t reflect.Type) ([]byte, error) {
	et := t.Elem()
	st := reflect.SliceOf(et)
	g.todo = append(g.todo, et, st)
	return execute(g.arrayTemplate, struct {
		Type, ElType, SliceType reflect.Type
		IsBytes                 bool
	}{
		Type:      t,
		ElType:    et,
		SliceType: st,
		IsBytes:   et == reflect.TypeOf(byte(0)),
	})
}

func (g *generator) genMap(t reflect.Type) ([]byte, error) {
	et := t.Elem()
	kt := t.Key()
	g.todo = append(g.todo, kt, et)
	return execute(g.mapTemplate, struct {
		Type, ElType, KeyType reflect.Type
	}{
		Type:    t,
		ElType:  et,
		KeyType: kt,
	})
}

func (g *generator) genMarshaler(t reflect.Type, kind string) ([]byte, error) {

	return execute(g.marshalTemplate, struct {
		Type reflect.Type
		Kind string
	}{
		Type: t,
		Kind: kind,
	})
}
func (g *generator) genStruct(t reflect.Type) ([]byte, error) {
	fn := g.typeIdentifier(t)
	var fields []field
	fields = exportedFields(t, g.fieldNames[fn])
	var names []string
	for _, f := range fields {
		names = append(names, f.Name)
		ft := f.Type
		if ft == nil {
			continue
		}
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		g.todo = append(g.todo, ft)
	}
	g.fieldNames[fn] = names // Update list field names.
	return execute(g.structTemplate, struct {
		Type, PtrType reflect.Type
		Fields        []field
	}{
		Type:    t,
		PtrType: reflect.PtrTo(t),
		Fields:  fields,
	})
}

// A field holds the information necessary to generate the encoder for a struct field.
type field struct {
	Name string
	Type reflect.Type
	Zero string // representation of the type's zero value
}

const fieldTagKey = "codec"

// exportedFields returns the exported fields of the struct type t that
// should be encoded, in the proper order.
// Exported fields of embedded, unexported types are not included.
// If there was a previous ordering, it is preserved, and new fields are
// added to the end.
// If a field was removed, we keep its number so as not to break existing
// encoded values. It will appear in the return value with an empty type.
//
// One drawback of this scheme is that it is not possible to rename a field.
// A rename will look like an addition and a removal.
func exportedFields(t reflect.Type, oldNames []string) []field {
	// Record the positions of the field names previously used for this struct,
	// so we can preserve them.
	fieldPos := map[string]int{}
	for i, n := range oldNames {
		fieldPos[n] = i
	}

	// If there are any new exported fields, assign them positions after the
	// existing ones.
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" { // Ignore unexported fields. A field is exported if its PkgPath is empty.
			continue
		}
		// Ignore a field if it has a struct tag with "-", like encoding/json.
		if tag, _ := f.Tag.Lookup(fieldTagKey); tag == "-" {
			continue
		}
		// Ignore fields of function and channel type.
		if f.Type.Kind() == reflect.Chan || f.Type.Kind() == reflect.Func {
			continue
		}
		if _, ok := fieldPos[f.Name]; !ok {
			fieldPos[f.Name] = len(fieldPos)
		}
	}

	// Populate the field structs, in the right order.
	fields := make([]field, len(fieldPos))
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if pos, ok := fieldPos[f.Name]; ok {
			fields[pos] = field{
				Name: f.Name,
				Type: f.Type,
				Zero: zeroValue(f.Type),
			}
		}
	}
	// Add back in the removed names, so their positions can be preserved in the
	// future.
	for i, n := range oldNames {
		if fields[i].Name == "" {
			fields[i].Name = n
		}
	}
	return fields
}

// zeroValue returns the string representation of a zero value of type t,
// or the empty string if there isn't one.
func zeroValue(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "false"
	case reflect.String:
		return `""`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "0"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "0"
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return "0"
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface:
		return "nil"
	default:
		return ""
	}
}

func execute(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encodeStatement returns a Go statement that encodes a value denoted by arg, of type t.
func (g *generator) encodeStatement(t reflect.Type, arg string) string {
	bn, native := builtinName(t)
	if bn != "" {
		// t can be handled by an Encoder method.
		if t != native {
			// t is not the Encoder method's argument type, so we must cast.
			arg = fmt.Sprintf("%s(%s)", native, arg)
		}
		return fmt.Sprintf("e.Encode%s(%s)", bn, arg)
	}
	if t.Kind() == reflect.Interface {
		return fmt.Sprintf("e.EncodeAny(%s)", arg)
	}
	if encodePtrArg(t) {
		arg = "&" + arg
	}
	if t.Name() != "" && t.Kind() != reflect.Struct && t.Kind() != reflect.Array {
		arg = fmt.Sprintf("%s(%s)", g.goName(t), arg)
	}
	return fmt.Sprintf("%s(e, %s)", g.encodeFunc(t), arg)
}

func encodePtrArg(t reflect.Type) bool {
	if t.Implements(binaryMarshalerType) || t.Implements(textMarshalerType) {
		return false
	}
	return t.Kind() == reflect.Struct || t.Kind() == reflect.Array
}

func (g *generator) encodeFunc(t reflect.Type) string {
	var typeName string
	bn, native := builtinName(t)
	if bn != "" {
		_ = native // TODO: ???
		typeName = "codecapi." + bn
	} else {
		typeName = g.typeIdentifier(t)
	}
	return fmt.Sprintf("(%s_codec{}).encode", typeName)
}

func (g *generator) decodeFunc(t reflect.Type) string {
	var typeName string
	bn, native := builtinName(t)
	if bn != "" {
		_ = native // TODO: ???
		typeName = "codecapi." + bn
	} else {
		typeName = g.typeIdentifier(t)
	}
	return fmt.Sprintf("(%s_codec{}).decode", typeName)
}

func (g *generator) decodeStatement(t reflect.Type, arg string) string {
	bn, native := builtinName(t)
	if bn != "" {
		// t can be handled by a Decoder method.
		if t != native {
			// t is not the Decoder method's return type, so we must cast.
			return fmt.Sprintf("%s = %s(d.Decode%s())", arg, g.goName(t), bn)
		}
		return fmt.Sprintf("%s = d.Decode%s()", arg, bn)
	}
	if t.Kind() == reflect.Interface {
		// t is an interface, so use DecodeAny, possibly with a type assertion.
		if t.NumMethod() == 0 {
			return fmt.Sprintf("%s = d.DecodeAny()", arg)
		}
		return fmt.Sprintf("%s = d.DecodeAny().(%s)", arg, g.goName(t))
	}
	// Assume we will generate a decode method for t.
	if t.Name() != "" && !willGenerate(t) {
		arg = fmt.Sprintf("(*%s)(&%s)", g.goName(t), arg)
	} else {
		arg = "&" + arg
	}
	return fmt.Sprintf("(%s_codec{}).decode(d, %s)", g.typeIdentifier(t), arg)
}

// willGenerate reports whether a codec will be generated for t.
func willGenerate(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		return true
	default:
		return false
	}
}

// builtinName returns the suffix to append to "encode" or "decode" to get the
// Encoder/Decoder method name for t. If t cannot be encoded by an Encoder
// method, the suffix is "". The second return value is the "native" type of the
// method: the argument to the Encoder method, and the return value of the
// Decoder method.
func builtinName(t reflect.Type) (suffix string, native reflect.Type) {
	switch t.Kind() {
	case reflect.String:
		return "String", reflect.TypeOf("")
	case reflect.Bool:
		return "Bool", reflect.TypeOf(true)
	case reflect.Int8, reflect.Uint8:
		return "Byte", reflect.TypeOf(byte(0))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return "Int", reflect.TypeOf(int64(0))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "Uint", reflect.TypeOf(uint64(0))
	case reflect.Float32, reflect.Float64:
		return "Float", reflect.TypeOf(0.0)
	case reflect.Slice:
		if t.Elem() == reflect.TypeOf(byte(0)) {
			return "Bytes", reflect.TypeOf([]byte(nil))
		}
	}
	return "", nil
}

// goName returns the name of t as it should appear in a Go program.
// E.g. "go/ast.File" => ast.File
// It assumes all package paths are represented in the file by their last element.
// TODO: relax that assumption.
func (g *generator) goName(t reflect.Type) string {
	if t.Name() != "" {
		s := t.String()
		if strings.HasPrefix(s, g.pkg+".") {
			s = s[len(g.pkg)+1:]
		}
		return s
	}
	switch t.Kind() {
	case reflect.Slice:
		return fmt.Sprintf("[]%s", g.goName(t.Elem()))
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), g.goName(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", g.goName(t.Key()), g.goName(t.Elem()))
	case reflect.Ptr:
		return fmt.Sprintf("*%s", g.goName(t.Elem()))
	default:
		panic(fmt.Sprintf("bad type: %s", t))
	}
}

var typeIdentifierReplacer = strings.NewReplacer("[]", "slice_", "[", "_", "]", "_", ".", "_", "*", "ptr_")

// typeIdentifier returns a valid Go identifier for type t.
// E.g. "ast.File" => "ast_File", "[]int" => "slice_int".
func (g *generator) typeIdentifier(t reflect.Type) string {
	n := typeIdentifierReplacer.Replace(g.goName(t))
	if t.Kind() == reflect.Array && t.Name() == "" {
		n = "array" + n
	}
	return n
}

// Template body for the beginning of the file.
const initialBody = `
// Code generated by the codec package. DO NOT EDIT.

package «.Package»

import (
	«range .Imports»
		"«.»"
	«- end»
)

`

// Template body for a slice type.
const sliceBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
type «$typeName» struct {}

func («$typeName») Init() {}

func (c «$typeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(«$goName»)) }

func (c «$typeName») encode(e *codecapi.Encoder, s «$goName») {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		«encodeStmt .ElType "x"»
	}
}

func (c «$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c «$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	n := d.StartList()
	if n < 0 { return }
	s := make([]«goName .ElType», n)
	for i := 0; i < n; i++ {
		«decodeStmt .ElType "s[i]"»
	}
	*p = s
}

func init() {
  codecapi.Register(«$goName»(nil), «$typeName»{})
}
`

// Template body for an array type.
const arrayBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
type «$typeName» struct {}

func («$typeName») Init() {}

func (c «$typeName») Encode(e *codecapi.Encoder, x interface{}) {
	a := x.(«$goName»)
	c.encode(e, &a)
}

func (c «$typeName») encode(e *codecapi.Encoder, s *«$goName») {
	«encodeStmt .SliceType "(*s)[:]"»
}

func (c «$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c «$typeName») decode(d *codecapi.Decoder, p *«$goName») {
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
			«decodeStmt .ElType "(*p)[i]"»
		}
	«end -»
}

func init() {
  codecapi.Register(«$goName»{}, «$typeName»{})
}
`

// Template body for a map type.
// A nil map is encoded as a zero.
// A map of size N is encoded as a list of length 2N, containing alternating
// keys and values.
//
// In the decode function, we declare a variable v to hold the decoded map value
// rather than decoding directly into m[v]. This is necessary for decode
// functions that take pointers: you can't take a pointer to a map element.
const mapBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
type «$typeName» struct{}

func (c «$typeName») Init() {}

func (c «$typeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(«$goName»)) }

func (c «$typeName») encode(e *codecapi.Encoder, m «$goName») {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2*len(m))
	for k, v := range m {
		«encodeStmt .KeyType "k"»
		«encodeStmt .ElType "v"»
	}
}

func (c «$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c «$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	n2 := d.StartList()
	if n2 < 0 { return }
	n := n2/2
	m := make(«$goName», n)
	var k «goName .KeyType»
	var v «goName .ElType»
	for i := 0; i < n; i++ {
		«decodeStmt .KeyType "k"»
		«decodeStmt .ElType "v"»
		m[k] = v
	}
	*p = m
}

func init() { codecapi.Register(«$goName»(nil), «$typeName»{}) }
`

// Template body for a type that implements encoding.BinaryMarshaler or encoding.TextMarshaler.
const marshalBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
type «$typeName» struct{}

func (c «$typeName») Init() {}

func (c «$typeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(«$goName»)) }

func (c «$typeName») encode(e *codecapi.Encoder, m «$goName») {
	data, err := m.Marshal«.Kind»()
	if err != nil {
		codecapi.Fail(err)
	}
	e.EncodeBytes(data)
}

func (c «$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c «$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	data := d.DecodeBytes()
	if err := p.Unmarshal«.Kind»(data); err != nil {
		codecapi.Fail(err)
	}
}

func init() { codecapi.Register(*new(«$goName»), «$typeName»{}) }
`

// Template body for a (pointer to a) struct type.
// A nil pointer is encoded as a zero. (This is done in Encoder.StartStruct.)
// Otherwise, a struct is encoded as the start code, its exported fields, then
// the end code. Each non-zero field is encoded as its field number followed by
// its value. A field that equals its zero value isn't encoded.
//
// The comment listing the field names is used when re-generating the file,
// to make sure we don't alter the existing mapping from field names to numbers.
const structBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $ptrTypeName := print "ptr_" $typeName »
« $goName := goName .Type »

// Fields of «$typeID»:«range .Fields» «.Name»«end»

type «$ptrTypeName» struct{}

func («$ptrTypeName») Init() {}

func (c «$ptrTypeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(*«$goName»)) }

func (c «$ptrTypeName») encode(e *codecapi.Encoder, x *«$goName») {
	if !e.StartPtr(x==nil, x) { return }
	(«$typeName»{}).encode(e, x)
}

func (c «$ptrTypeName») Decode(d *codecapi.Decoder) interface{} {
	var x *«$goName»
	c.decode(d, &x)
	return x
}

func (c «$ptrTypeName») decode(d *codecapi.Decoder, p **«$goName») {
	proceed, ref := d.StartPtr()
	if !proceed { return }
	if ref != nil {
		*p = ref.(*«$goName»)
		return
	}
	var x «$goName»
	d.StoreRef(&x)
	(«$typeName»{}).decode(d, &x)
	*p = &x
}

type «$typeName» struct{}

func («$typeName») Init() {}

func (c «$typeName») Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(«$goName»)
	c.encode(e, &s)
 }

func (c «$typeName») encode(e *codecapi.Encoder, x *«$goName») {
	e.StartStruct()
	«range $i, $f := .Fields»
		«- if $f.Type -»
			«- if $f.Zero -»
				if x.«$f.Name» != «$f.Zero» {
			«- end»
			e.EncodeUint(«$i»)
			«encodeStmt .Type (print "x." $f.Name)»
			«- if $f.Zero -»
			}
			«- end»
		«- end»
	«end -»
	e.EndStruct()
}

func (c «$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c «$typeName») decode(d *codecapi.Decoder, x *«$goName») {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 { break }
		switch n {
		«range $i, $f := .Fields -»
			«- if $f.Type -»
				case «$i»:
		        «decodeStmt $f.Type (print "x." $f.Name)»
			«end -»
		«end -»
		default:
			d.UnknownField("«$goName»", n)
		}
	}
}



func init() {
	codecapi.Register(«$goName»{}, «$typeName»{})
	codecapi.Register(&«$goName»{}, «$ptrTypeName»{})
}
`
