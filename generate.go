// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/jba/codec/codecapi"
)

//go:generate ./embed.sh

type GenerateOptions struct {
	// FieldTag is the name that GenerateFile will use to look up
	// field tag information. The default is "codec".
	FieldTag string
}

// GenerateFile writes encoders and decoders to filename. It generates code for
// the type of each given value, as well as any types they depend on.
// packagePath is the output package path.
//
// The encoding assigns numbers to struct fields for efficient decoding.
// GenerateFile reads filename if it exists to discover the field numbers that
// have already been assigned, so it can preserve them. So it is important that
// the existing output file remains available to the generator, or changes to
// your structs may result in existing encoded data being decoded incorrectly.
func GenerateFile(filename, packagePath string, opts *GenerateOptions, values ...interface{}) error {
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	fieldTag := "codec"
	if opts != nil && opts.FieldTag != "" {
		fieldTag = opts.FieldTag
	}
	if err := generate(f, packagePath, fieldTag, values...); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func generate(w io.Writer, packagePath string, fieldTag string, vs ...interface{}) error {
	g := &generator{
		pkgPath:     packagePath,
		fieldTagKey: fieldTag,
	}
	funcs := template.FuncMap{
		"typeID":     g.typeID,
		"goName":     g.goName,
		"encodeStmt": g.encodeStmt,
		"decodeStmt": g.decodeStmt,
		"encodeFunc": g.encodeFunc,
	}

	newTemplate := func(name, body string) *template.Template {
		return template.Must(template.New(name).Delims("«", "»").Funcs(funcs).Parse(body))
	}

	g.initialTemplate = newTemplate("initial", initialBody)
	g.sliceTemplate = newTemplate("slice", sliceBody)
	g.arrayTemplate = newTemplate("array", arrayBody)
	g.mapTemplate = newTemplate("map", mapBody)
	g.ptrTemplate = newTemplate("ptr", ptrBody)
	g.structTemplate = newTemplate("struct", structBody)
	g.marshalTemplate = newTemplate("marshaler", marshalBody)

	src, err := g.generate(vs)
	if err != nil {
		return err
	}
	fsrc, err := format.Source(src)
	if err != nil {
		filename, err2 := writeTempFile("bad-source-*.go", src)
		var msg string
		if err2 != nil {
			msg = fmt.Sprintf("could not write bad source: %v", err)
		} else {
			msg = fmt.Sprintf("wrote bad source to %s", filename)
		}
		return fmt.Errorf("format.Source: %v;\n%s", err, msg)
	}
	_, err = w.Write(fsrc)
	return err
}

func writeTempFile(pattern string, contents []byte) (string, error) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(contents); err != nil {
		return "", err
	}
	return f.Name(), nil
}

type generator struct {
	pkgPath         string
	fieldTagKey     string
	importMap       map[string]string // import path to import identifier
	pkgPathMap      map[string]string //package path to qualifying identifier
	initialTemplate *template.Template
	sliceTemplate   *template.Template
	arrayTemplate   *template.Template
	mapTemplate     *template.Template
	ptrTemplate     *template.Template
	structTemplate  *template.Template
	marshalTemplate *template.Template
}

type importSpec struct {
	Path, ID string
}

func (g *generator) generate(typevals []interface{}) ([]byte, error) {
	todo := g.referencedTypeList(typevals)
	g.buildImportMap(todo)
	var code []byte
	for _, t := range todo {
		piece, err := g.gen(t)
		if err != nil {
			return nil, err
		}
		if piece != nil {
			header := fmt.Sprintf("//// %s\n\n", t)
			code = append(code, header...)
			code = append(code, piece...)
		}
	}

	var stdImports, otherImports []importSpec
	for path, id := range g.importMap {
		if path == g.pkgPath {
			continue
		}
		spec := importSpec{path, id}
		if strings.ContainsRune(path, '.') {
			otherImports = append(otherImports, spec)
		} else {
			stdImports = append(stdImports, spec)
		}
	}
	initial, err := execute(g.initialTemplate, struct {
		Package                  string
		StdImports, OtherImports []importSpec
	}{
		Package:      path.Base(g.pkgPath),
		StdImports:   stdImports,
		OtherImports: otherImports,
	})
	if err != nil {
		return nil, err
	}
	return append(initial, code...), nil
}

// referencedTypeList returns a list of all types referenced from typevals.
func (g *generator) referencedTypeList(typevals []interface{}) []reflect.Type {
	// Collect all the types referred to, except builtins. We will generate most
	// of these (not defined types whose underlying type is builtin, for
	// example), but we need them all to generate the right import statements.
	types := map[reflect.Type]bool{}
	for _, v := range typevals {
		g.referencedTypes(reflect.TypeOf(v), types)
	}
	var typeList []reflect.Type
	for t := range types {
		typeList = append(typeList, t)
	}
	// Sort for determinism.
	sort.Slice(typeList, func(i, j int) bool {
		return codecapi.TypeString(typeList[i], nil) < codecapi.TypeString(typeList[j], nil)
	})
	return typeList
}

// referencedTypes records in the set m all the types referenced from t.
func (g *generator) referencedTypes(t reflect.Type, m map[reflect.Type]bool) {
	if m[t] {
		return
	}
	switch t.Kind() {
	case reflect.Slice:
		if t.Name() == "" && t.Elem() == byteType {
			return
		}
		m[t] = true
		g.referencedTypes(t.Elem(), m)
	case reflect.Ptr:
		m[t] = true
		g.referencedTypes(t.Elem(), m)
	case reflect.Array:
		m[t] = true
		g.referencedTypes(t.Elem(), m)
		g.referencedTypes(reflect.SliceOf(t.Elem()), m)
	case reflect.Map:
		m[t] = true
		g.referencedTypes(t.Key(), m)
		g.referencedTypes(t.Elem(), m)
	case reflect.Struct:
		m[t] = true
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !g.ignoreField(t, f) {
				g.referencedTypes(f.Type, m)
			}
		}
	default:
		if t.PkgPath() != "" {
			m[t] = true
		}
	}
}

func packageName(t reflect.Type) string {
	if t.PkgPath() == "" {
		return ""
	}
	s := t.String()
	i := strings.LastIndexByte(s, '.')
	if i < 0 {
		panic(fmt.Sprintf("type %s has non-empty PkgPath but no dot in String", t))
	}
	return s[:i]
}

func (g *generator) ignoreField(structType reflect.Type, f reflect.StructField) bool {
	// Ignore unexported fields for structs in a different package. A field
	// is exported if its PkgPath is empty.
	if structType.PkgPath() != g.pkgPath && f.PkgPath != "" {
		return true
	}
	// Ignore fields of function and channel type.
	if f.Type.Kind() == reflect.Chan || f.Type.Kind() == reflect.Func {
		return true
	}
	// Ignore a field if it has a struct tag with "-", like encoding/json.
	_, omit := parseTag(g.fieldTagKey, f.Tag)
	return omit
}

func (g *generator) buildImportMap(types []reflect.Type) {
	g.importMap = map[string]string{
		"reflect":                       "",
		"github.com/jba/codec/codecapi": "",
	}
	// Collect the prefixes in use so far.
	// For these, assume that the package names are the last components of the
	// import paths.
	prefixes := map[string]bool{}
	for ppath, id := range g.importMap {
		if id == "" {
			prefixes[path.Base(ppath)] = true
		} else {
			prefixes[id] = true
		}
	}
	for _, t := range types {
		ppath := t.PkgPath()
		if ppath == "" {
			continue
		}
		if ppath == g.pkgPath {
			continue
		}
		if _, ok := g.importMap[ppath]; ok {
			continue
		}
		// Determine an import identifier for the path.
		var id string
		// The package prefix used in the file will be the package name, unless
		// we provide an import identifier. Usually, the package name is the
		// last component of the import path.
		prefix := path.Base(ppath)
		// For package names that differ from their path's last component,
		// provide the name as an import identifier, to simplify code
		// generation.
		pname := packageName(t)
		if pname != prefix {
			prefix = pname
			id = pname
		}
		// If the prefix is not unique, generate a unique one for the
		// identifier.
		orig := prefix
		for i := 1; prefixes[prefix]; i++ {
			prefix = fmt.Sprintf("%s%d", orig, i)
			id = prefix
		}
		prefixes[prefix] = true
		g.importMap[ppath] = id
	}
	// The package path map is used to generate Go names for types. It is close
	// to the import map, but not the same: first, it includes g.pkgPath.
	// Second, a package mapping to an empty string in the import map,
	// indicating no import identifier, will map to its last component in the
	// path map.
	g.pkgPathMap = map[string]string{g.pkgPath: ""}
	for k, v := range g.importMap {
		if v == "" {
			v = path.Base(k)
		}
		g.pkgPathMap[k] = v
	}
}

var (
	binaryMarshalerType   = reflect.TypeOf(new(encoding.BinaryMarshaler)).Elem()
	binaryUnmarshalerType = reflect.TypeOf(new(encoding.BinaryUnmarshaler)).Elem()
	textMarshalerType     = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
	textUnmarshalerType   = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
	byteType              = reflect.TypeOf(byte(0))
)

func (g *generator) gen(t reflect.Type) ([]byte, error) {
	if m := implementsMarshaler(t); m != "" {
		return g.genMarshaler(t, m)
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
		return g.genPtr(t)
	}
	return nil, nil
}

// willGenerate reports whether a codec will be generated for t.
func willGenerate(t reflect.Type) bool {
	if implementsMarshaler(t) != "" {
		return true
	}
	switch t.Kind() {
	case reflect.Slice:
		return t.Elem() != byteType
	case reflect.Struct, reflect.Array, reflect.Map, reflect.Ptr:
		return true
	default:
		return false
	}
}

// implementsMarshaler returns the kind of Marshaler that t implements ("Binary"
// or "Text"), or the empty string if it doesn't implement one.
func implementsMarshaler(t reflect.Type) string {
	if t.Implements(binaryMarshalerType) && reflect.PtrTo(t).Implements(binaryUnmarshalerType) {
		return "Binary"
	}
	if t.Implements(textMarshalerType) && reflect.PtrTo(t).Implements(textUnmarshalerType) {
		return "Text"
	}
	return ""
}

func (g *generator) genSlice(t reflect.Type) ([]byte, error) {
	return execute(g.sliceTemplate, struct {
		Type    reflect.Type
		ElField bool
	}{
		Type:    t,
		ElField: willGenerate(t.Elem()),
	})
}

func (g *generator) genArray(t reflect.Type) ([]byte, error) {
	et := t.Elem()
	st := reflect.SliceOf(et)
	return execute(g.arrayTemplate, struct {
		Type, SliceType reflect.Type
		IsBytes         bool
		ElField         bool
	}{
		Type:      t,
		SliceType: st,
		IsBytes:   et == byteType,
		ElField:   willGenerate(et),
	})
}

func (g *generator) genMap(t reflect.Type) ([]byte, error) {
	et := t.Elem()
	kt := t.Key()
	return execute(g.mapTemplate, struct {
		Type              reflect.Type
		KeyField, ElField bool
	}{
		Type:     t,
		KeyField: willGenerate(kt),
		ElField:  willGenerate(et) && kt != et,
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

func (g *generator) genPtr(t reflect.Type) ([]byte, error) {
	return execute(g.ptrTemplate, struct {
		Type    reflect.Type
		ElField bool
	}{
		Type:    t,
		ElField: willGenerate(t.Elem()),
	})
}

func (g *generator) genStruct(t reflect.Type) ([]byte, error) {
	if t.Name() == "" {
		return nil, fmt.Errorf("cannot generate code for unnamed struct type %s", t)
	}
	fields := g.structFields(t)
	fieldTypesSet := map[reflect.Type]bool{}
	for _, f := range fields {
		ft := f.Type
		if ft == nil {
			continue
		}
		if willGenerate(ft) {
			fieldTypesSet[ft] = true
		}
	}
	var fieldTypes []reflect.Type
	for t := range fieldTypesSet {
		fieldTypes = append(fieldTypes, t)
	}
	// Sort so the list is deterministic, for testing. The strings returned by
	// reflect.Type.String aren't unique (e.g. []pkg.Foo where there are two
	// packages with name "pkg"), but that doesn't matter as long as no tests
	// trigger the problem.
	sort.Slice(fieldTypes, func(i, j int) bool {
		return fieldTypes[i].String() < fieldTypes[j].String()
	})
	return execute(g.structTemplate, struct {
		Type, PtrType reflect.Type
		Fields        []field
		FieldTypes    []reflect.Type // unique list of types
	}{
		Type:       t,
		PtrType:    reflect.PtrTo(t),
		Fields:     fields,
		FieldTypes: fieldTypes,
	})
}

// A field holds the information necessary to generate the encoder for a struct field.
// This struct's fields are exported so they can be used in templates.
type field struct {
	Name string
	Type reflect.Type
	Zero string // representation of the type's zero value
}

// structFields returns the fields of the struct type t that should be encoded.
// For structs in a package other than the one being generated into, that
// includes all direct exported fields, but not exported fields of embedded,
// unexported types. For structs in the same package, unexported fields are
// included.
func (g *generator) structFields(t reflect.Type) []field {
	var fields []field
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if g.ignoreField(t, f) {
			continue
		}
		name, _ := parseTag(g.fieldTagKey, f.Tag)
		if name == "" {
			name = f.Name
		}
		fields = append(fields, field{
			Name: name,
			Type: f.Type,
			Zero: zeroValue(f.Type),
		})
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

// encodeStmt returns a Go statement that encodes a value denoted by arg, of type t.
func (g *generator) encodeStmt(t reflect.Type, arg string) string {
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
	// If the encode function expects a pointer, take the address of the arg.
	if encodePtrArg(t) {
		if arg[0] == '*' {
			// If the arg is a dereference, just remove the dereference.
			arg = arg[1:]
		} else {
			arg = "&" + arg
		}
	}
	return fmt.Sprintf("%s(e, %s)", g.encodeFunc(t), arg)
}

// encodePtrArg reports whether the type is passed by pointer.
// We pass potentially large values by pointer for efficiency.
func encodePtrArg(t reflect.Type) bool {
	if t.Implements(binaryMarshalerType) || t.Implements(textMarshalerType) {
		return false
	}
	return t.Kind() == reflect.Struct || t.Kind() == reflect.Array
}

func (g *generator) encodeFunc(t reflect.Type) string {
	var typeName string
	bn, _ := builtinName(t)
	if bn != "" {
		typeName = "codecapi." + bn
	} else {
		typeName = g.typeID(t)
	}
	return fmt.Sprintf("c.%s_codec.encode", typeName)
}

func (g *generator) decodeStmt(t reflect.Type, arg string) string {
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
	return fmt.Sprintf("c.%s_codec.decode(d, %s)", g.typeID(t), arg)
}

// builtinName returns the suffix to append to "encode" or "decode" to get the
// Encoder/Decoder method name for t. If t cannot be encoded by an Encoder
// method, the suffix is "". The second return value is the "native" type of the
// method: the argument to the Encoder method, and the return value of the
// Decoder method.
func builtinName(t reflect.Type) (suffix string, native reflect.Type) {
	if implementsMarshaler(t) != "" {
		return "", nil
	}
	switch t.Kind() {
	case reflect.String:
		return "String", reflect.TypeOf("")
	case reflect.Bool:
		return "Bool", reflect.TypeOf(true)
	case reflect.Int8, reflect.Uint8:
		return "Byte", byteType
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return "Int", reflect.TypeOf(int64(0))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "Uint", reflect.TypeOf(uint64(0))
	case reflect.Float32, reflect.Float64:
		return "Float", reflect.TypeOf(0.0)
	case reflect.Complex64, reflect.Complex128:
		return "Complex", reflect.TypeOf(0 + 0i)
	case reflect.Slice:
		if t.Elem() == byteType {
			return "Bytes", reflect.TypeOf([]byte(nil))
		}
	}
	return "", nil
}

// goName returns the name of t as it should appear in a Go program.
// E.g. "go/ast.File" => ast.File
func (g *generator) goName(t reflect.Type) string {
	return codecapi.TypeString(t, g.pkgPathMap)
}

var typeIDReplacer = strings.NewReplacer(
	"[]", "slice_",
	"{}", "", // for empty interface
	"[", "_", "]", "_", ".", "_",
	"*", "ptr_",
)

// typeID returns a valid Go identifier for type t.
// E.g. "ast.File" => "ast_File", "[]int" => "slice_int".
func (g *generator) typeID(t reflect.Type) string {
	if t.Name() != "" {
		return strings.ReplaceAll(g.goName(t), ".", "_")
	}
	switch t.Kind() {
	case reflect.Slice:
		return "slice_" + g.typeID(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("array_%d_%s", t.Len(), g.typeID(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map_%s__%s", g.typeID(t.Key()), g.typeID(t.Elem()))
	case reflect.Ptr:
		return "ptr_" + g.typeID(t.Elem())
	default:
		return typeIDReplacer.Replace(g.goName(t))
	}
}

// parseTag extracts the sub-tag named by key, then parses it using the
// de facto standard format introduced in encoding/json:
//   "-" means "ignore this tag". It must occur by itself. (parseTag returns an error
//       in this case, whereas encoding/json accepts the "-" even if it is not alone.)
//   "<name>" provides an alternative name for the field
//   "<name>,opt1,opt2,..." specifies options after the name.
// The return values are:
// name: the name given in tag, or "" if there is no name.
// omit: true if the field should be omitted.
// options: the list of options.
func parseTag(key string, t reflect.StructTag) (name string, omit bool) {
	s := t.Get(key)
	parts := strings.Split(s, ",")
	if parts[0] == "-" {
		return "", true
	}
	return parts[0], false
}
