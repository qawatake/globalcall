package globalcall

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/qawatake/globalcall/internal/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const name = "globalcall"
const doc = "globalcall detects that specific functions are called in a package scope"
const url = "https://pkg.go.dev/github.com/qawatake/globalcall"

func NewAnalyzer(funcs ...Func) *analysis.Analyzer {
	r := runner{
		funcs: funcs,
	}
	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		URL:  url,
		Run:  r.run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}
}

type runner struct {
	funcs []Func
}

type Func struct {
	// Package path of the function or method.
	PkgPath string
	// Name of the function or method.
	FuncName string
}

func (r *runner) run(pass *analysis.Pass) (any, error) {

	m := make(map[types.Object]Func)
	for _, fn := range r.funcs {
		obj, err := funcObjectOf(pass, fn)
		if err != nil {
			if errors.Is(err, errFuncNotUsed) {
				continue
			}
			return nil, err
		}
		m[obj] = fn
	}
	// ↓の各行を集めたものがvalues
	// var (
	// x = hoge()
	// y = 2
	// a, b = fuga(), 3
	// )
	values := valueSpecs(pass.Files)
	for _, value := range values {
		for _, n := range value.Values { // var x, y = hoge(), fuga()の右辺たち
			if call := methodCall(n); call != nil {
				if f, ok := m[pass.TypesInfo.ObjectOf(call)]; ok {
					reportf(pass, call.Pos(), f.FuncName)
				}
				continue
			}
			if call := funcCall(n); call != nil {
				if f, ok := m[pass.TypesInfo.ObjectOf(call)]; ok {
					reportf(pass, call.Pos(), f.FuncName)
				}
				continue
			}
		}
	}
	return nil, nil
}

func reportf(pass *analysis.Pass, pos token.Pos, funcName string) {
	pass.Reportf(pos, "%s must not be called in a package scope", funcName)
}

func methodCall(expr ast.Expr) (ident *ast.Ident) {
	defer func() {
		if e := recover(); e != nil {
			ident = nil
		}
	}()
	return expr.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel
}

func funcCall(expr ast.Expr) (ident *ast.Ident) {
	defer func() {
		if e := recover(); e != nil {
			ident = nil
		}
	}()
	return expr.(*ast.CallExpr).Fun.(*ast.Ident)
}

func valueSpecs(fs []*ast.File) []*ast.ValueSpec {
	v := make([][]*ast.ValueSpec, 0, len(fs))
	for _, f := range fs {
		vv := make([][]*ast.ValueSpec, 0, len(f.Decls))
		for _, dec := range f.Decls {
			vv = append(vv, values(dec))
		}
		v = append(v, slices.Concat(vv...))
	}
	return slices.Concat(v...)
}

func values(dec ast.Decl) (v []*ast.ValueSpec) {
	defer func() {
		if e := recover(); e != nil {
			v = nil
		}
	}()
	specs := dec.(*ast.GenDecl).Specs
	vs := make([]*ast.ValueSpec, 0, len(specs))
	for _, spec := range specs {
		if v, ok := spec.(*ast.ValueSpec); ok {
			vs = append(vs, v)
		}
	}
	return vs
}

func funcObjectOf(pass *analysis.Pass, d Func) (*types.Func, error) {
	// function
	if !strings.Contains(d.FuncName, ".") {
		obj := analysisutil.ObjectOf(pass, d.PkgPath, d.FuncName)
		if obj == nil {
			// not found is ok because func need not to be called.
			return nil, errFuncNotUsed
		}
		ft, ok := obj.(*types.Func)
		if !ok {
			return nil, newErrNotFunc(d.PkgPath, d.FuncName)
		}
		return ft, nil
	}
	tt := strings.Split(d.FuncName, ".")
	if len(tt) != 2 {
		return nil, newErrInvalidFuncName(d.FuncName)
	}
	// method
	recv := tt[0]
	method := tt[1]
	recvType := analysisutil.TypeOf(pass, d.PkgPath, recv)
	if recvType == nil {
		return nil, errFuncNotUsed
	}
	m := analysisutil.MethodOf(recvType, method)
	if m == nil {
		return nil, errFuncNotUsed
	}
	return m, nil
}

var errFuncNotUsed = errors.New("function not used")

type errInvalidFuncName struct {
	FuncName string
}

func newErrInvalidFuncName(funcName string) errInvalidFuncName {
	return errInvalidFuncName{
		FuncName: funcName,
	}
}

func (e errInvalidFuncName) Error() string {
	return fmt.Sprintf("invalid FuncName %s", e.FuncName)
}

type errNotFunc struct {
	PkgPath  string
	FuncName string
}

func newErrNotFunc(pkgPath, funcName string) errNotFunc {
	return errNotFunc{
		PkgPath:  pkgPath,
		FuncName: funcName,
	}
}

func (e errNotFunc) Error() string {
	return fmt.Sprintf("%s.%s is not a function.", e.PkgPath, e.FuncName)
}
