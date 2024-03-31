package optimizer

import "github.com/uchijo/bf-eval/instr"

func Optimize(src []instr.Instruction) []instr.Instruction {
	src = ResetToZeroPattern(src)
	src = SumShift(src)
	return src
}