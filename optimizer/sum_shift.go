package optimizer

import "github.com/uchijo/bf-eval/instr"

// SumShift returns a new slice of instr.Instruction that sums up the shift operations.
func SumShift(src []instr.Instruction) []instr.Instruction {
	retval := []instr.Instruction{}
	for i := 0; i < len(src); i++ {
		if src[i].Op == instr.OpShiftRight {
			j := SearchCont(src, instr.OpShiftRight, i)
			retval = append(retval, instr.Instruction{
				Op:   instr.OpShiftRight,
				Data: uint8(j - i),
			})
			i = j - 1
		} else if src[i].Op == instr.OpShiftLeft {
			j := SearchCont(src, instr.OpShiftLeft, i)
			retval = append(retval, instr.Instruction{
				Op:   instr.OpShiftLeft,
				Data: uint8(j - i),
			})
			i = j - 1
		} else {
			retval = append(retval, src[i])
		}
	}
	return retval
}

// SearchCont returns the index of the next instruction that is not the same as the given op.
func SearchCont(src []instr.Instruction, op instr.Op, start int) int {
	for i := start; i < len(src); i++ {
		if src[i].Op != op {
			return i
		}
	}
	return len(src)
}
