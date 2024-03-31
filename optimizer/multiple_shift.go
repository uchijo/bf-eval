package optimizer

import "github.com/uchijo/bf-eval/instr"

// 以下のパターンを探して、OpMultiShiftに割り付ける
// LoopStart
//   ShiftLeft 9
// LoopEnd

func MultipleShift(src []instr.Instruction) []instr.Instruction {
	retval := []instr.Instruction{}
	for i := 0; i < len(src); i++ {
		if matchMultipleShift(src, i) {
			var shift int16
			if src[i+1].Op == instr.OpShiftRight {
				shift = int16(src[i+1].Data)
			} else {
				shift = -int16(src[i+1].Data)
			}
			retval = append(retval, instr.Instruction{
				Op:   instr.OpMultiShift,
				Data: shift,
			})
			i += 2
		} else {
			retval = append(retval, src[i])
		}
	}
	return retval
}

func matchMultipleShift(src []instr.Instruction, i int) bool {
	if i+2 >= len(src) {
		return false
	}
	return src[i].Op == instr.OpLoopStart &&
		(src[i+1].Op == instr.OpShiftLeft || src[i+1].Op == instr.OpShiftRight) &&
		src[i+2].Op == instr.OpLoopEnd
}
