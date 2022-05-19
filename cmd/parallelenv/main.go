package main

import (
	"github.com/sho-hata/parallelenv"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(parallelenv.Analyzer) }
