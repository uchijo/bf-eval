package optimizer

import "github.com/uchijo/bf-eval/instr"

// do same as SumShift but for + and -
func SumIncrDecr(src []instr.Instruction) []instr.Instruction {
	retval := []instr.Instruction{}
	for i := 0; i < len(src); i++ {
		if src[i].Op == instr.OpIncr {
			j := SearchCont(src, instr.OpIncr, i)
			retval = append(retval, instr.Instruction{
				Op:   instr.OpIncr,
				Data: uint8(j - i),
			})
			i = j - 1
		} else if src[i].Op == instr.OpDecr {
			j := SearchCont(src, instr.OpDecr, i)
			retval = append(retval, instr.Instruction{
				Op:   instr.OpDecr,
				Data: uint8(j - i),
			})
			i = j - 1
		} else {
			retval = append(retval, src[i])
		}
	}
	return retval
}
