package globalcall_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/qawatake/globalcall"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, globalcall.NewAnalyzer(
		globalcall.Func{
			PkgPath:  "a",
			FuncName: "call",
		},
		globalcall.Func{
			PkgPath:  "a",
			FuncName: "X.call",
		},
		globalcall.Func{
			PkgPath:  "a",
			FuncName: "*X.Call",
		},
	), "a")
}
