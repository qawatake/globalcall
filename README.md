# globalcall

[![Go Reference](https://pkg.go.dev/badge/github.com/qawatake/globalcall.svg)](https://pkg.go.dev/github.com/qawatake/globalcall)
[![test](https://github.com/qawatake/globalcall/actions/workflows/test.yaml/badge.svg)](https://github.com/qawatake/globalcall/actions/workflows/test.yaml)

Linter `globalcall` detects that specific functions are called in a package scope.

```go
var i = Int() // ng because Int must not be called in a package scope.

func main() {
  j := Int() // ok
  fmt.Println(j)
}
```

## How to use

Build your `globalcall` binary by writing `main.go` like below.

```go
package main

import (
  "github.com/qawatake/globalcall"
  "golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
  unitchecker.Main(
    globalcall.NewAnalyzer(
      globalcall.Func{
        PkgPath:  "pkg/in/which/target/func/is/defined",
        FuncName: "Call",
      },
    ),
  )
}
```

Then, run `go vet` with your `globalcall` binary.

```sh
go vet -vettool=/path/to/your/globalcall ./...
```
