package main

import (
	"github.com/qawatake/globalcall"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(globalcall.NewAnalyzer(
		globalcall.Func{
			PkgPath:  "github.com/qawatake/globalcall/internal/example_test",
			FuncName: "Call",
		},
		globalcall.Func{
			PkgPath:  "github.com/qawatake/globalcall/internal/example_test",
			FuncName: "X.Call",
		},
		globalcall.Func{
			PkgPath:  "github.com/qawatake/globalcall/internal/example_test",
			FuncName: "*Y.Call",
		},
	))
}
