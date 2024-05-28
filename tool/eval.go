package tool

import (
	"os"

	"github.com/uchijo/bf-eval/instr"
	"github.com/uchijo/bf-eval/optimizer"
)

func Eval(src []instr.Instruction) {
	buf := []byte{}
	w := os.Stdout

	mem := NewMemStore()
	var memPtr int32 = 0
	pc := 0
	src = optimizer.Optimize(src)
	jumpDest := cacheJumpDest(src)
	programLen := len(src)

	for {
		if pc >= programLen {
			break
		}

		op := src[pc].Op
		if op == instr.OpLoopEnd {
			if mem.Get(memPtr) != 0 {
				pc = jumpDest[pc]
			}
			goto done
		}
		if op == instr.OpShiftRight {
			memPtr += src[pc].Data
			goto done
		}
		if op == instr.OpShiftLeft {
			memPtr -= src[pc].Data
			goto done
		}
		if op == instr.OpAddMem {
			mem.AddTo(memPtr+src[pc].Data, mem.Get(memPtr))
			mem.Set(memPtr, 0)
			goto done
		}
		if op == instr.OpLoopStart {
			if mem.Get(memPtr) == 0 {
				pc = jumpDest[pc]
			}
			goto done
		}
		if op == instr.OpIncr {
			mem.AddTo(memPtr, uint8(src[pc].Data))
			goto done
		}
		if op == instr.OpDecr {
			mem.SubFrom(memPtr, uint8(src[pc].Data))
			goto done
		}
		if op == instr.OpZeroReset {
			mem.Set(memPtr, 0)
			goto done
		}
		if op == instr.OpMultiShift {
			for mem.Get(memPtr) != 0 {
				memPtr += src[pc].Data
			}
			goto done
		}
		if op == instr.OpSubMem {
			mem.SubFrom(memPtr+src[pc].Data, mem.Get(memPtr))
			mem.Set(memPtr, 0)
			goto done
		}
		if op == instr.OpOutput {
			buf = append(buf, mem.Get(memPtr))
			if len(buf) >= 4096 {
				w.Write(buf)
				buf = []byte{}
			}
			goto done
		}
		
		done:
		pc++
	}
	if len(buf) > 0 {
		w.Write(buf)
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
