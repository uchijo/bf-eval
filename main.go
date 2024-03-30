package main

import (
	"fmt"
	"os"
)

func main() {
	// read from file
	b, err := os.ReadFile("m.bf")
	if err != nil {
		panic(err)
	}

	eval(b)
}

func eval(src []uint8) {
	mem := map[int]uint8{}
	memPtr := 0
	pc := 0

	for {
		if pc >= len(src) {
			break
		}

		switch src[pc] {
		case '>':
			memPtr++
		case '<':
			memPtr--
		case '+':
			mem[memPtr]++
		case '-':
			mem[memPtr]--
		case '.':
			fmt.Print(string(mem[memPtr]))
		case ',':
			// not implemented
		case '[':
			if mem[memPtr] == 0 {
				// find matching ']'
				nest := 1
				for {
					pc++
					if src[pc] == '[' {
						nest++
					} else if src[pc] == ']' {
						nest--
						if nest == 0 {
							break
						}
					}
				}
			}
		case ']':
			if mem[memPtr] != 0 {
				// find matching '['
				nest := 1
				for {
					pc--
					if src[pc] == ']' {
						nest++
					} else if src[pc] == '[' {
						nest--
						if nest == 0 {
							break
						}
					}
				}
			}
		}

		pc++
	}
}
