package tool

import (
	"bufio"
	"fmt"
	"os"

	"github.com/uchijo/bf-eval/instr"
)

func Eval(src []instr.Instruction) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	mem := map[int]uint8{}
	memPtr := 0
	pc := 0
	src = resetToZeroPattern(src)
	jumpDest := cacheJumpDest(src)

	for {
		if pc >= len(src) {
			break
		}

		switch src[pc].Op {
		case instr.OpShiftRight:
			memPtr++
		case instr.OpShiftLeft:
			memPtr--
		case instr.OpIncr:
			mem[memPtr]++
		case instr.OpDecr:
			mem[memPtr]--
		case instr.OpOutput:
			fmt.Fprint(w, string(mem[memPtr]))
		case instr.OpInput:
			// not implemented
		case instr.OpZeroReset:
			mem[memPtr] = 0
		case instr.OpLoopStart:
			if mem[memPtr] == 0 {
				pc = jumpDest[pc]
			}
		case instr.OpLoopEnd:
			if mem[memPtr] != 0 {
				pc = jumpDest[pc]
			}
		}

		pc++
	}
}

// find [-] pattern and replace it with 0
func resetToZeroPattern(src []instr.Instruction) []instr.Instruction {
	retval := []instr.Instruction{}
	for i := 0; i < len(src); i++ {
		if i+2 < len(src) && src[i].Op == instr.OpLoopStart && src[i+1].Op == instr.OpDecr && src[i+2].Op == instr.OpLoopEnd {
			retval = append(retval, instr.Instruction{Op: instr.OpZeroReset})
			i += 2
		} else {
			retval = append(retval, src[i])
		}
	}
	return retval
}

func cacheJumpDest(src []instr.Instruction) map[int]int {
	jumpDest := map[int]int{}
	for pc, c := range src {
		if c.Op == instr.OpLoopStart {
			start := pc
			nest := 1
			for {
				start++
				if src[start].Op == instr.OpLoopStart {
					nest++
				} else if src[start].Op == instr.OpLoopEnd {
					nest--
					if nest == 0 {
						break
					}
				} else {
					continue
				}
			}
			jumpDest[pc] = start
		} else if c.Op == instr.OpLoopEnd {
			start := pc
			nest := 1
			for {
				start--
				if src[start].Op == instr.OpLoopEnd {
					nest++
				} else if src[start].Op == instr.OpLoopStart {
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
