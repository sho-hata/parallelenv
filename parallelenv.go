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
				report(pass, fnDecl)
			}
		}
	})

	return nil, nil
}

type target int

const (
	unknown target = iota
	setEnv
	parallel
)

func (t target) isValid() bool {
	switch t {
	case setEnv, parallel:
		return true
	case unknown:
		return false
	}
	return false
}

func sToTarget(s string) target {
	if t, ok := map[string]target{"t.Setenv": setEnv, "t.Parallel": parallel}[s]; ok {
		return t
	}

	return unknown
}

// checkState
// manages the location in the source code where the check target appears.
type checkState map[target][]token.Pos

func newCheckState() checkState {
	return checkState{setEnv: nil, parallel: nil}
}

// shouldReport
// Returns whether to report it as an error or not
func (s checkState) shouldReport() bool {
	if pos, ok := s[setEnv]; !ok || len(pos) == 0 {
		return false
	}

	if pos, ok := s[parallel]; !ok || len(pos) == 0 {
		return false
	}

	return true
}

// report
// In the Body part of the function declaration, check whether t.Parallel and t.SetEnv appear at the same time.
func report(pass *analysis.Pass, fnDecl *ast.FuncDecl) {
	state := newCheckState()
	for _, stmt := range fnDecl.Body.List {
		state = checkStmt(stmt, state)
	}
	if state.shouldReport() {
		for _, pos := range state {
			for _, p := range pos {
				pass.Reportf(p, "cannot set environment variables in parallel tests")
			}
		}
	}
}

func checkStmt(stmt ast.Stmt, state checkState) checkState {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		state = shouldReportExprStmt(stmt, state)
	case *ast.IfStmt:
		state = shouldReportIfStmt(stmt, state)
	case *ast.BlockStmt:
		for _, s := range stmt.List {
			state = checkStmt(s, state)
		}
	}
	return state
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

	target := sToTarget(ident.Name + "." + fn.Sel.Name)
	if ok := target.isValid(); ok {
		checkState[target] = append(checkState[target], stmt.Pos())
	}

	return checkState
}

func shouldReportIfStmt(stmt *ast.IfStmt, state checkState) checkState {
	for _, s := range stmt.Body.List {
		state = checkStmt(s, state)
	}

	if (stmt.Else) != nil {
		state = checkStmt(stmt.Else, state)
	}

	return state
}
