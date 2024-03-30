package main

import (
	"os"

	"github.com/uchijo/bf-eval/tool"
)

func main() {
	// read from file
	b, err := os.ReadFile("m.bf")
	if err != nil {
		panic(err)
	}
	parsed := tool.Parse(b)

	tool.Eval(parsed)
}
