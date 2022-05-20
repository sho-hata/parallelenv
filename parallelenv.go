package parallelenv

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "parallelenv is go linter that check whether environment variables are set in tests run in parallel"

var Analyzer = &analysis.Analyzer{
	Name: "parallelenv",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			if !strings.HasSuffix(pass.Fset.File(n.Pos()).Name(), "_test.go") {
				return
			}

			for _, decl := range n.Decls {
				fnDecl, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				checkFnDecl(pass, fnDecl)
			}
		}
	})

	return nil, nil
}

type checkState struct {
	isSetEnv   bool
	isParallel bool
	positions  []token.Pos
}

func (s *checkState) shouldReport() bool {
	return s.isSetEnv && s.isParallel
}

// checkFnDecl
// In the Body part of the function declaration, check whether t.Parallel and t.SetEnv appear at the same time.
func checkFnDecl(pass *analysis.Pass, fnDecl *ast.FuncDecl) {
	state := checkState{}
	for _, stmt := range fnDecl.Body.List {
		switch stmt := stmt.(type) {
		case *ast.ExprStmt:
			state = shouldReportExprStmt(stmt, state)
		}
	}
	if state.shouldReport() {
		for _, pos := range state.positions {
			pass.Reportf(pos, "cannot set environment variables in parallel tests")
		}
	}
}

func shouldReportExprStmt(stmt *ast.ExprStmt, checkState checkState) checkState {
	callExpr, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return checkState
	}
	fn, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return checkState
	}

	ident, ok := fn.X.(*ast.Ident)
	if !ok {
		return checkState
	}
	targetName := ident.Name + "." + fn.Sel.Name

	if targetName == "t.Parallel" {
		checkState.isParallel = true
		checkState.positions = append(checkState.positions, stmt.Pos())
	}
	if targetName == "t.Setenv" {
		checkState.isSetEnv = true
		checkState.positions = append(checkState.positions, stmt.Pos())
	}

	return checkState
}
