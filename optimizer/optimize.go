package optimizer

import "github.com/uchijo/bf-eval/instr"

func Optimize(src []instr.Instruction) []instr.Instruction {
	src = ResetToZeroPattern(src)
	src = SumShift(src)
	src = SumIncrDecr(src)
	src = FindAddMem(src)
	src = MultipleShift(src)
	return src
}
