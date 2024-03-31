package tool

import "github.com/uchijo/bf-eval/instr"

func DumpInstr(src []instr.Instruction) {
	for _, i := range src {
		switch i.Op {
		case instr.OpShiftRight:
			println("ShiftRight", i.Data)
		case instr.OpShiftLeft:
			println("ShiftLeft", i.Data)
		case instr.OpIncr:
			println("Incr")
		case instr.OpDecr:
			println("Decr")
		case instr.OpOutput:
			println("Output")
		case instr.OpInput:
			println("Input")
		case instr.OpZeroReset:
			println("ZeroReset")
		case instr.OpLoopStart:
			println("LoopStart")
		case instr.OpLoopEnd:
			println("LoopEnd")
		}
	}
}
