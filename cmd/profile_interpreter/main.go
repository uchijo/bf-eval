package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

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

	// write result to ../result-<timestamp>.pprof
	filename := fmt.Sprintf("profile/result-%d.pprof", time.Now().Unix())
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// eval
	pprof.StartCPUProfile(f)
	tool.Eval(parsed)
	pprof.StopCPUProfile()
}
