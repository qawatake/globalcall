package main

import (
	"github.com/qawatake/globalcall"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(globalcall.Analyzer) }
