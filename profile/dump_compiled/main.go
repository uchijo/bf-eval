package main

import (
	"os"

	"github.com/uchijo/bf-eval/optimizer"
	"github.com/uchijo/bf-eval/tool"
)

func main() {
	// read from file
	b, err := os.ReadFile("m.bf")
	if err != nil {
		panic(err)
	}
	parsed, err := tool.Parse(b)
	if err != nil {
		panic(err)
	}

	optimized := optimizer.Optimize(parsed)
	tool.DumpInstr(optimized)
}
