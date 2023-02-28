package main

import (
	"fmt"
	"github.com/panpito/momock/generator"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"strings"
)

func main() {
	pwd, _ := os.Getwd()
	fileName := os.Getenv("GOFILE")

	filePath := fmt.Sprint(pwd, "/", fileName)
	log.Printf("generating mock for: %s", filePath)

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	a := momock_generator.Ast{
		PackageName: astFile.Name.Name,
		Imports:     astFile.Imports,
		Interfaces:  []*ast.TypeSpec{},
	}

	ast.Inspect(astFile, func(n ast.Node) bool {
		switch t := n.(type) {
		// find variable declarations
		case *ast.TypeSpec:
			// which are public
			if t.Name.IsExported() {
				switch t.Type.(type) {
				// and are interfaces
				case *ast.InterfaceType:
					a.Interfaces = append(a.Interfaces, t)
				}
			}
		}

		return true
	})

	aMock := a.ToMock()
	for _, spec := range a.Imports {
		astutil.AddImport(fset, &aMock, strings.Trim(spec.Path.Value, "\""))
	}
	astutil.AddImport(fset, &aMock, "testing")
	astutil.AddImport(fset, &aMock, "github.com/panpito/momock/manager")
	astutil.AddImport(fset, &aMock, "log")
	ast.SortImports(fset, &aMock)

	mockFile := strings.ReplaceAll(filePath, ".go", "_mock.go")
	f, _ := os.Create(mockFile)
	printer.Fprint(f, fset, &aMock)
}
