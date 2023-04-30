package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
)

func main() {
	fs := token.NewFileSet()
	astFile, err := parser.ParseFile(fs, "proto/opq_host.pb.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	astutil.Apply(astFile, nil, func(cursor *astutil.Cursor) bool {
		n := cursor.Node()
		switch v := n.(type) {
		case *ast.FuncDecl:
			if v.Recv != nil {
				if v.Name.Name == "Load" {
					v.Name.Name = "LoadWithBytes"
					v.Type.Params.List[1].Names[0].Name = "b"
					v.Type.Params.List[1].Type.(*ast.Ident).Name = "[]byte"
					v.Body.List = v.Body.List[2:]
				}
			}
		case *ast.GenDecl:
			if v.Tok.String() == "import" {
				for i, j := range v.Specs {
					imp := j.(*ast.ImportSpec)
					if imp.Name.Name == "os" {
						v.Specs = append(v.Specs[:i], v.Specs[i+1:]...)
					}
				}
			}

		}
		return true
	})
	f, _ := os.OpenFile("proto/opq_host.pb.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	defer f.Close()
	err = format.Node(f, fs, astFile)
	if err != nil {
		panic(err)
	}
}
