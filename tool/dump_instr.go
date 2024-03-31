package tool

import (
	"fmt"

	"github.com/uchijo/bf-eval/instr"
)

func DumpInstr(src []instr.Instruction) {
	for _, i := range src {
		switch i.Op {
		case instr.OpShiftRight:
			fmt.Println("ShiftRight", i.Data)
		case instr.OpShiftLeft:
			fmt.Println("ShiftLeft", i.Data)
		case instr.OpIncr:
			fmt.Println("Incr")
		case instr.OpDecr:
			fmt.Println("Decr")
		case instr.OpOutput:
			fmt.Println("Output")
		case instr.OpInput:
			fmt.Println("Input")
		case instr.OpZeroReset:
			fmt.Println("ZeroReset")
		case instr.OpLoopStart:
			fmt.Println("LoopStart")
		case instr.OpLoopEnd:
			fmt.Println("LoopEnd")
		}
	}
}
