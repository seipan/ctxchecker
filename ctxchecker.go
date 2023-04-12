package ctxchecker

import (
	"go/ast"
	"strconv"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "ctxchecker is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "ctxchecker",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	_ = getCommentMap(pass)
	pkgs := pass.Pkg.Imports()
	obj := analysisutil.LookupFromImports(pkgs, "context", "Context")
	types := pass.TypesInfo

	inspect.Preorder(nil, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FieldList:
			for _, v := range n.List {
				value, ok := v.Type.(*ast.SelectorExpr)
				if !ok {
					return
				}
				if types.ObjectOf(value.Sel) == obj {
					pass.Reportf(n.Pos(), "ctx is here")
				}

			}
		}
	})

	return nil, nil
}

func getCommentMap(pass *analysis.Pass) map[string]string {
	var mp = make(map[string]string)

	for _, file := range pass.Files {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				pos := pass.Fset.Position(c.Pos())
				mp[pos.Filename+"_"+strconv.Itoa(pos.Line)] = c.Text
			}
		}
	}
	return mp
}
