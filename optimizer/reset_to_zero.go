package optimizer

import "github.com/uchijo/bf-eval/instr"

// find [-] pattern and replace it with 0
func ResetToZeroPattern(src []instr.Instruction) []instr.Instruction {
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
