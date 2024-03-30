package tool

import "fmt"

func Eval(src []uint8) {
	mem := map[int]uint8{}
	memPtr := 0
	pc := 0
	src = resetToZeroPattern(src)
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
		case '0':
			mem[memPtr] = 0
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

// find [-] pattern and replace it with 0
func resetToZeroPattern(src []uint8) []uint8 {
	retval := []uint8{}
	for i := 0; i < len(src); i++ {
		if i+2 < len(src) && src[i] == '[' && src[i+1] == '-' && src[i+2] == ']' {
			retval = append(retval, '0')
			i += 2
		} else {
			retval = append(retval, src[i])
		}
	}
	return retval
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
