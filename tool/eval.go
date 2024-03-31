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

	mem := NewMemStore()
	memPtr := 0
	pc := 0
	src = optimizer.Optimize(src)
	jumpDest := cacheJumpDest(src)
	programLen := len(src)

	for {
		if pc >= programLen {
			break
		}

		switch src[pc].Op {
		case instr.OpLoopEnd:
			if mem.Get(memPtr) != 0 {
				pc = jumpDest[pc]
			}
		case instr.OpShiftRight:
			memPtr += int(src[pc].Data)
		case instr.OpShiftLeft:
			memPtr -= int(src[pc].Data)
		case instr.OpCopy:
			mem.AddTo(memPtr+int(src[pc].Data), mem.Get(memPtr))
			mem.Set(memPtr, 0)
		case instr.OpLoopStart:
			if mem.Get(memPtr) == 0 {
				pc = jumpDest[pc]
			}
		case instr.OpIncr:
			mem.AddTo(memPtr, uint8(src[pc].Data))
		case instr.OpDecr:
			mem.SubFrom(memPtr, uint8(src[pc].Data))
		case instr.OpZeroReset:
			mem.Set(memPtr, 0)
		case instr.OpMultiShift:
			for mem.Get(memPtr) != 0 {
				memPtr += int(src[pc].Data)
			}
		case instr.OpOutput:
			fmt.Fprint(w, string(mem.Get(memPtr)))
			// case instr.OpInput:
			// 	counts[instr.OpInput]++
			// 	// not implemented
		}

		pc++
	}
}

func cacheJumpDest(src []instr.Instruction) []int {
	jumpDest := make([]int, len(src))
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
