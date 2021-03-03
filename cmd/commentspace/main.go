package main

import (
	"github.com/pmatseykanets/go-checks/checks/commentspace"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(commentspace.Analyzer)
}
