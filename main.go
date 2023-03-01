package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/panpito/momock/generator"
	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("could not get current path: %v", err)
	}
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
	f, err := os.Create(mockFile)
	if err != nil {
		log.Printf("could not create mock file: %v", err)
	}

	var b bytes.Buffer
	if err := printer.Fprint(bufio.NewWriter(&b), fset, &aMock); err != nil {
		log.Printf("could not write mock: %v", err)
	}

	formattedFileBytes, err := format.Source(b.Bytes())
	if err != nil {
		log.Printf("could not format: %v", err)
	}
	if _, err := f.Write(formattedFileBytes); err != nil {
		log.Printf("could not write formatted file: %v", err)
	}
}
