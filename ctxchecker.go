package ctxchecker

import (
	"go/ast"
	"reflect"
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

	inspect.Preorder(nil, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FieldList:
			if !checkHandlerNoPoint(pass, n) && !checkTest(pass, n) && !ctxCheck(pass, n) && !checkHandler(pass, n) {
				pass.Reportf(n.Pos(), "no ctx")
			}

		}
	})

	return nil, nil
}

func ctxCheck(pass *analysis.Pass, field *ast.FieldList) bool {
	flag := false
	pkgs := pass.Pkg.Imports()
	Obj := analysisutil.LookupFromImports(pkgs, "context", "Context")
	types := pass.TypesInfo
	for _, v := range field.List {
		value, ok := v.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if reflect.DeepEqual(types.ObjectOf(value.Sel), Obj) {
			flag = true
		}
	}
	return flag
}

func checkHandlerNoPoint(pass *analysis.Pass, field *ast.FieldList) bool {
	pkgs := pass.Pkg.Imports()
	httpObj := analysisutil.LookupFromImports(pkgs, "net/http", "Request")
	ginObj := analysisutil.LookupFromImports(pkgs, "gin", "Context")
	types := pass.TypesInfo
	for _, v := range field.List {
		svalue, ok := v.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		value, ok := svalue.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if types.ObjectOf(value.Sel) == httpObj || types.ObjectOf(value.Sel) == ginObj {
			return true
		}
	}

	return false
}

func checkHandler(pass *analysis.Pass, field *ast.FieldList) bool {
	pkgs := pass.Pkg.Imports()
	httpObj := analysisutil.LookupFromImports(pkgs, "net/http", "ResponseWriter")
	echoObj := analysisutil.LookupFromImports(pkgs, "echo", "Context")
	types := pass.TypesInfo
	for _, v := range field.List {
		value, ok := v.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if types.ObjectOf(value.Sel) == httpObj || types.ObjectOf(value.Sel) == echoObj {
			return true
		}
	}

	return false
}

func checkTest(pass *analysis.Pass, field *ast.FieldList) bool {
	flag := false
	pkgs := pass.Pkg.Imports()
	Obj := analysisutil.LookupFromImports(pkgs, "testing", "T")
	types := pass.TypesInfo
	for _, v := range field.List {
		value, ok := v.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		selvalue, ok := value.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if reflect.DeepEqual(types.ObjectOf(selvalue.Sel), Obj) {
			flag = true
		}
	}
	return flag
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
