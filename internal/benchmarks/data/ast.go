// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

// Go parse trees: the ast.Files from the net/http package.
// These have cycles, which we remove before saving.

import (
	"encoding/gob"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var ASTData = gobBenchmarkData("ast", func() interface{} { return new(map[string]*ast.File) })

var astTypes = []interface{}{
	ast.ArrayType{},
	ast.AssignStmt{},
	ast.BadDecl{},
	ast.BadExpr{},
	ast.BadStmt{},
	ast.BasicLit{},
	ast.BinaryExpr{},
	ast.BlockStmt{},
	ast.BranchStmt{},
	ast.CallExpr{},
	ast.CaseClause{},
	ast.ChanType{},
	ast.CommClause{},
	ast.CommentGroup{},
	ast.Comment{},
	ast.CompositeLit{},
	ast.DeclStmt{},
	ast.DeferStmt{},
	ast.Ellipsis{},
	ast.EmptyStmt{},
	ast.ExprStmt{},
	ast.FieldList{},
	ast.Field{},
	ast.ForStmt{},
	ast.FuncDecl{},
	ast.FuncLit{},
	ast.FuncType{},
	ast.GenDecl{},
	ast.GoStmt{},
	ast.Ident{},
	ast.IfStmt{},
	ast.ImportSpec{},
	ast.IncDecStmt{},
	ast.IndexExpr{},
	ast.InterfaceType{},
	ast.KeyValueExpr{},
	ast.LabeledStmt{},
	ast.MapType{},
	ast.ParenExpr{},
	ast.RangeStmt{},
	ast.ReturnStmt{},
	ast.Scope{},
	ast.SelectStmt{},
	ast.SelectorExpr{},
	ast.SendStmt{},
	ast.SliceExpr{},
	ast.StarExpr{},
	ast.StructType{},
	ast.SwitchStmt{},
	ast.TypeAssertExpr{},
	ast.TypeSpec{},
	ast.TypeSwitchStmt{},
	ast.UnaryExpr{},
	ast.ValueSpec{},
}

func init() {
	for _, n := range []interface{}{
		&ast.ArrayType{},
		&ast.AssignStmt{},
		&ast.BadDecl{},
		&ast.BadExpr{},
		&ast.BadStmt{},
		&ast.BasicLit{},
		&ast.BinaryExpr{},
		&ast.BlockStmt{},
		&ast.BranchStmt{},
		&ast.CallExpr{},
		&ast.CaseClause{},
		&ast.ChanType{},
		&ast.CommClause{},
		&ast.CommentGroup{},
		&ast.Comment{},
		&ast.CompositeLit{},
		&ast.DeclStmt{},
		&ast.DeferStmt{},
		&ast.Ellipsis{},
		&ast.EmptyStmt{},
		&ast.ExprStmt{},
		&ast.FieldList{},
		&ast.Field{},
		&ast.ForStmt{},
		&ast.FuncDecl{},
		&ast.FuncLit{},
		&ast.FuncType{},
		&ast.GenDecl{},
		&ast.GoStmt{},
		&ast.Ident{},
		&ast.IfStmt{},
		&ast.ImportSpec{},
		&ast.IncDecStmt{},
		&ast.IndexExpr{},
		&ast.InterfaceType{},
		&ast.KeyValueExpr{},
		&ast.LabeledStmt{},
		&ast.MapType{},
		&ast.ParenExpr{},
		&ast.RangeStmt{},
		&ast.ReturnStmt{},
		&ast.Scope{},
		&ast.SelectStmt{},
		&ast.SelectorExpr{},
		&ast.SendStmt{},
		&ast.SliceExpr{},
		&ast.StarExpr{},
		&ast.StructType{},
		&ast.SwitchStmt{},
		&ast.TypeAssertExpr{},
		&ast.TypeSpec{},
		&ast.TypeSwitchStmt{},
		&ast.UnaryExpr{},
		&ast.ValueSpec{},
	} {
		gob.Register(n)
	}
}

// ParseStdlibPackage parses and returns the standard library package at ppath.
// It assumes the package name is the last component of the path.
func ParseStdlibPackage(ppath string) (*ast.Package, error) {
	dir := filepath.Join(runtime.GOROOT(), "src", "net", "http")
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return pkgs[path.Base(ppath)], nil
}

func generateASTToFile(filename string) error {
	// Get the AST for the net/http package.
	httpPkg, err := ParseStdlibPackage("net/http")
	if err != nil {
		return err
	}
	// nil out things that result in cycles.
	for _, f := range httpPkg.Files {
		f.Scope = nil
		ast.Inspect(f, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				id.Obj = nil
			}
			return true
		})
	}
	return WriteNewFile(filename, func(f *os.File) error {
		return gob.NewEncoder(f).Encode(httpPkg.Files)
	})
}
