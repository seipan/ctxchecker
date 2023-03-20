package main

import (
	"github.com/seipan/ctxchecker"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(ctxchecker.Analyzer) }
