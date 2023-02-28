package momock_generator

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Ast struct {
	PackageName string
	Imports     []*ast.ImportSpec
	Interfaces  []*ast.TypeSpec
}

func (a *Ast) ToMock() ast.File {
	declarations := make([]ast.Decl, 0)
	for _, interfaceType := range a.Interfaces {
		declarations = append(declarations, mockInterface(interfaceType)...)
	}

	return ast.File{
		Package: token.NoPos,
		Name:    ast.NewIdent(a.PackageName),
		Decls:   declarations,
		Scope:   nil,
	}
}

func mockInterface(interfaceTypeSpec *ast.TypeSpec) []ast.Decl {
	mockStructName := ast.NewIdent(fmt.Sprint("Mock", interfaceTypeSpec.Name.Name))
	mockConstructorName := ast.NewIdent(fmt.Sprint("New", mockStructName))
	xMockLib := ast.NewIdent("momock_manager")
	selMockLibCallerData := ast.NewIdent("CallerData")
	selMockManager := ast.NewIdent("MockManager")
	selNewMockManager := ast.NewIdent("NewMockManager")

	mockStruct := &ast.GenDecl{
		TokPos: token.NoPos,
		Tok:    token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: mockStructName,
				Type: &ast.StructType{
					Struct: token.NoPos,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("t")},
								Type: &ast.StarExpr{
									Star: token.NoPos,
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("testing"),
										Sel: ast.NewIdent("T"),
									},
								},
							},
							{
								Names: []*ast.Ident{selMockManager},
								Type: &ast.StarExpr{
									Star: token.NoPos,
									X: &ast.SelectorExpr{
										X:   xMockLib,
										Sel: selMockManager,
									},
								},
							},
						},
						Opening: token.NoPos,
						Closing: token.NoPos,
					},
				},
			},
		},
	}

	mockConstructor := &ast.FuncDecl{
		Name: mockConstructorName,
		Type: &ast.FuncType{
			Func: token.NoPos,
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("t")},
						Type: &ast.StarExpr{
							Star: token.NoPos,
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("testing"),
								Sel: ast.NewIdent("T"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							Star: token.NoPos,
							X:    mockStructName,
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: mockStructName,
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("t"),
										Value: ast.NewIdent("t"),
									},
									&ast.KeyValueExpr{
										Key: selMockManager,
										Value: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   xMockLib,
												Sel: selNewMockManager,
											},
											Args: []ast.Expr{
												ast.NewIdent("t"),
											},
										},
									},
								},
								Incomplete: false,
							},
						},
					},
				},
			},
		},
	}

	receiverName := ast.NewIdent("mock")
	receiverFieldMockManager := ast.NewIdent("MockManager")
	receiverField := &ast.FieldList{
		Opening: token.NoPos,
		Closing: token.NoPos,
		List: []*ast.Field{
			{
				Names: []*ast.Ident{receiverName},
				Type: &ast.StarExpr{
					Star: token.NoPos,
					X:    mockStructName,
				},
			},
		},
	}

	i := interfaceTypeSpec.Type.(*ast.InterfaceType)

	var mockFunctions []ast.Decl
	for _, method := range i.Methods.List {
		switch method.Type.(type) {
		case *ast.FuncType:
			castedFunction := method.Type.(*ast.FuncType)
			inputs := describeFieldList(castedFunction.Params)
			outputs := describeFieldList(castedFunction.Results)

			callerDataAssignment := &ast.AssignStmt{
				TokPos: token.NoPos,
				Tok:    token.DEFINE,
				Lhs:    []ast.Expr{ast.NewIdent("callerData")},
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   xMockLib,
							Sel: selMockLibCallerData,
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: ast.NewIdent("MethodName"),
								Value: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   receiverName,
											Sel: receiverFieldMockManager,
										},
										Sel: ast.NewIdent("WhatsMyName"),
									},
								},
							},
							&ast.KeyValueExpr{
								Key:   ast.NewIdent("InputsLength"),
								Value: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(inputs.length)},
							},
							&ast.KeyValueExpr{
								Key:   ast.NewIdent("OutputsLength"),
								Value: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(outputs.length)},
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Inputs"),
								Value: &ast.CompositeLit{
									Type: &ast.MapType{
										Map:   token.NoPos,
										Key:   ast.NewIdent("int"),
										Value: ast.NewIdent("any"),
									},
									Elts: inputs.inputAsMap,
								},
							},
						},
					},
				},
			}

			verifyAssignment := &ast.AssignStmt{
				TokPos: token.NoPos,
				Tok:    token.DEFINE,
				Lhs:    []ast.Expr{ast.NewIdent("out")},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   receiverName,
								Sel: receiverFieldMockManager,
							},
							Sel: ast.NewIdent("Verify"),
						},
						Args: []ast.Expr{ast.NewIdent("callerData")},
					},
				},
			}

			logVerifyAssignment := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("log"),
						Sel: ast.NewIdent("Print"),
					},
					Args: []ast.Expr{ast.NewIdent("out")},
				},
			}

			returnStatment := &ast.ReturnStmt{
				Results: outputs.outputArgs,
			}

			var statementsList []ast.Stmt
			statementsList = append([]ast.Stmt{callerDataAssignment, verifyAssignment, logVerifyAssignment}, outputs.outputAsAssignment...)
			statementsList = append(statementsList, returnStatment)

			funcDecl := &ast.FuncDecl{
				Recv: receiverField,
				Name: method.Names[0],
				Type: method.Type.(*ast.FuncType),
				Body: &ast.BlockStmt{
					List: statementsList,
				},
			}

			mockFunctions = append(mockFunctions, funcDecl)
		}
	}

	return append([]ast.Decl{mockStruct, mockConstructor}, mockFunctions...)
}

type simplifiedFieldList struct {
	length             int
	inputAsMap         []ast.Expr
	outputArgs         []ast.Expr
	outputAsAssignment []ast.Stmt
}

func describeFieldList(fieldList *ast.FieldList) (l simplifiedFieldList) {
	if fieldList == nil {
		return l
	}

	list := fieldList.List
	l.length = len(list)

	if len(list) != 0 {
		listAsExpr := make([]ast.Expr, len(list))
		outputArgs := make([]ast.Expr, len(list))
		outputAsAssignment := make([]ast.Stmt, len(list))

		for i, field := range list {
			var value string
			if len(field.Names) != 0 {
				value = field.Names[0].Name
			} else {
				value = "nope"
			}

			listAsExpr[i] = &ast.KeyValueExpr{
				Key: &ast.BasicLit{
					Kind:  token.INT,
					Value: fmt.Sprint(i),
				},
				Value: ast.NewIdent(value),
			}

			returnArg := ast.NewIdent(fmt.Sprint("return", i))
			outputArgs[i] = returnArg

			outputAsAssignment[i] = &ast.AssignStmt{
				TokPos: token.NoPos,
				Tok:    token.DEFINE,
				Lhs: []ast.Expr{
					returnArg,
					ast.NewIdent("_"),
				},
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X: &ast.IndexExpr{
							X: ast.NewIdent("out"),
							Index: &ast.BasicLit{
								Kind:  token.INT,
								Value: fmt.Sprint(i),
							},
						},
						Type: field.Type,
					},
				},
			}
		}

		l.inputAsMap = listAsExpr
		l.outputArgs = outputArgs
		l.outputAsAssignment = outputAsAssignment
	}

	return l
}
