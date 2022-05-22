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
				if fnDecl, ok := decl.(*ast.FuncDecl); ok {
					analyze(pass, fnDecl)
				}
			}
		}
	})

	return nil, nil
}

type target int

const (
	unknown target = iota
	tSetEnv
	tParallel
)

func (t target) isValid() bool {
	switch t {
	case tSetEnv, tParallel:
		return true
	case unknown:
		return false
	}

	return false
}

func sToTarget(s string) target {
	if t, ok := map[string]target{"t.Setenv": tSetEnv, "t.Parallel": tParallel}[s]; ok {
		return t
	}

	return unknown
}

type testLevel = int

const mainTest testLevel = 0

// results
// manages the location in the source code where the check target appears.
type results map[testLevel][]*result

type result struct {
	target target
	pos    token.Pos
}

func (rs results) ShouldReportResults() results {
	for testLevel, res := range rs {
		var callSetEnv, callParallel int

		for _, r := range res {
			switch r.target {
			case tSetEnv:
				callSetEnv++
			case tParallel:
				callParallel++
			case unknown:
			}
		}

		if callSetEnv == 0 || callParallel == 0 {
			delete(rs, testLevel)
		}
	}

	return rs
}

func (rs results) report(pass *analysis.Pass) {
	for _, results := range rs {
		for _, result := range results {
			pass.Reportf(result.pos, "cannot set environment variables in parallel tests")
		}
	}
}

// analyze
// In the Body part of the function declaration, check whether t.Parallel and t.SetEnv appear at the same time.
func analyze(pass *analysis.Pass, fnDecl *ast.FuncDecl) {
	state := results{}

	for _, stmt := range fnDecl.Body.List {
		state = state.check(stmt, mainTest)
	}

	state.ShouldReportResults().report(pass)
}

func (rs results) check(stmt ast.Stmt, testLevel testLevel) results {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		rs = rs.checkExprStmt(stmt, testLevel)
	case *ast.IfStmt:
		rs = rs.checkIfStmt(stmt, testLevel)
	case *ast.BlockStmt:
		for _, s := range stmt.List {
			rs = rs.check(s, testLevel)
		}
	case *ast.RangeStmt:
		rs = rs.check(stmt.Body, testLevel)
	}

	return rs
}

func (rs results) checkExprStmt(stmt *ast.ExprStmt, testLevel testLevel) results {
	callExpr, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return rs
	}

	fn, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return rs
	}

	ident, ok := fn.X.(*ast.Ident)
	if !ok {
		return rs
	}

	fnName := ident.Name + "." + fn.Sel.Name

	if fnName == "t.Run" {
		testLevel++

		for _, expr := range callExpr.Args {
			switch e := expr.(type) {
			case *ast.FuncLit:
				for _, s := range e.Body.List {
					rs = rs.check(s, testLevel)
				}
			}
		}

		return rs
	}

	target := sToTarget(fnName)
	if ok := target.isValid(); ok {
		rs[testLevel] = append(rs[testLevel], &result{target: target, pos: stmt.Pos()})
	}

	return rs
}

func (rs results) checkIfStmt(stmt *ast.IfStmt, testLevel testLevel) results {
	for _, s := range stmt.Body.List {
		rs = rs.check(s, testLevel)
	}

	if (stmt.Else) != nil {
		rs = rs.check(stmt.Else, testLevel)
	}

	return rs
}
