package tool

import "fmt"

func Eval(src []uint8) {
	mem := map[int]uint8{}
	memPtr := 0
	pc := 0
	jumpDest := cacheJumpDest(src)

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
				pc = jumpDest[pc]
			}
		case ']':
			if mem[memPtr] != 0 {
				pc = jumpDest[pc]
			}
		}

		pc++
	}
}

func cacheJumpDest(src []uint8) map[int]int {
	jumpDest := map[int]int{}
	for pc, c := range src {
		if c == '[' {
			start := pc
			nest := 1
			for {
				start++
				if src[start] == '[' {
					nest++
				} else if src[start] == ']' {
					nest--
					if nest == 0 {
						break
					}
				} else {
					continue
				}
			}
			jumpDest[pc] = start
		} else if c == ']' {
			start := pc
			nest := 1
			for {
				start--
				if src[start] == ']' {
					nest++
				} else if src[start] == '[' {
					nest--
					if nest == 0 {
						break
					}
				} else {
					continue
				}
			}
			jumpDest[pc] = start
		}
	}
	return jumpDest
}
