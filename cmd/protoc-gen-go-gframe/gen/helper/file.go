package helper

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ParseGoFile(filepath string) ([]byte, map[string]struct{}) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil
	}
	gosrc, err := io.ReadAll(file)
	if err != nil {
		return nil, nil
	}
	fnMap := make(map[string]struct{})
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", gosrc, 0)
	if err != nil {
		return gosrc, nil
	}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			fnMap[fn.Name.String()] = struct{}{}
		}
	}
	return gosrc, fnMap
}
