package optimizer

import (
	"github.com/uchijo/bf-eval/instr"
)

// 以下のパターンは、現在のポインタ位置の値をn個右のセルに足す操作と等価
// LoopStart
// 	Decr 1
// 	ShiftRight n
// 	Incr 1
// 	ShiftLeft n
// LoopEnd
// これをOpAddMemに割り付ける

func FindAddMem(src []instr.Instruction) []instr.Instruction {
	var copy []instr.Instruction
	for i := 0; i < len(src); i++ {
		if matchAddMem(src, i) {
			var newInst instr.Instruction
			if src[i+2].Op == instr.OpShiftRight {
				newInst = instr.Instruction{Op: instr.OpAddMem, Data: src[i+2].Data}
			} else {
				newInst = instr.Instruction{Op: instr.OpAddMem, Data: -src[i+2].Data}
			}
			copy = append(copy, newInst)
			i += 5
		} else {
			copy = append(copy, src[i])
		}
	}
	return copy
}

func matchAddMem(src []instr.Instruction, pos int) bool {
	if pos+5 >= len(src) {
		return false
	}

	// 命令列がパターンに一致するか確認
	rMatches := src[pos].Op == instr.OpLoopStart &&
		src[pos+1].Op == instr.OpDecr &&
		src[pos+2].Op == instr.OpShiftRight &&
		src[pos+3].Op == instr.OpIncr &&
		src[pos+4].Op == instr.OpShiftLeft &&
		src[pos+5].Op == instr.OpLoopEnd
	lMatches := src[pos].Op == instr.OpLoopStart &&
		src[pos+1].Op == instr.OpDecr &&
		src[pos+2].Op == instr.OpShiftLeft &&
		src[pos+3].Op == instr.OpIncr &&
		src[pos+4].Op == instr.OpShiftRight &&
		src[pos+5].Op == instr.OpLoopEnd
	// どちらかのパターンに一致しない場合はfalse
	if !(rMatches || lMatches) {
		return false
	}

	// インクリメント数が適正か確認
	if src[pos+1].Data != 1 || src[pos+3].Data != 1 {
		return false
	}

	// シフト数が適正か
	return src[pos+2].Data == src[pos+4].Data
}
