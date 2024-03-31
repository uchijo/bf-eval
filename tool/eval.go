package tool

import (
	"bufio"
	"fmt"
	"os"

	"github.com/uchijo/bf-eval/instr"
	"github.com/uchijo/bf-eval/optimizer"
)

func Eval(src []instr.Instruction) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	mem := map[int]uint8{}
	memPtr := 0
	pc := 0
	src = optimizer.Optimize(src)
	jumpDest := cacheJumpDest(src)

	for {
		if pc >= len(src) {
			break
		}

		switch src[pc].Op {
		case instr.OpShiftRight:
			memPtr += int(src[pc].Data)
		case instr.OpShiftLeft:
			memPtr -= int(src[pc].Data)
		case instr.OpIncr:
			mem[memPtr] += uint8(src[pc].Data)
		case instr.OpDecr:
			mem[memPtr] -= uint8(src[pc].Data)
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
		case instr.OpCopy:
			mem[memPtr+int(src[pc].Data)] += mem[memPtr]
			mem[memPtr] = 0
		}

		pc++
	}
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
