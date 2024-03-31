package instr

import (
	"errors"
	"fmt"
)

type Op uint8

const (
	OpLoopEnd Op = iota
	OpShiftRight
	OpShiftLeft
	OpAddMem
	OpLoopStart
	OpIncr
	OpDecr
	OpZeroReset
	OpMultiShift
	OpOutput
	OpInput
)

type Instruction struct {
	Op   Op
	pad  uint8
	Data int16
}

func NewInstruction(ch uint8) (Instruction, error) {
	switch ch {
	case '>':
		return Instruction{Op: OpShiftRight, Data: 1}, nil
	case '<':
		return Instruction{Op: OpShiftLeft, Data: 1}, nil
	case '+':
		return Instruction{Op: OpIncr, Data: 1}, nil
	case '-':
		return Instruction{Op: OpDecr, Data: 1}, nil
	case '.':
		return Instruction{Op: OpOutput}, nil
	case ',':
		return Instruction{Op: OpInput}, nil
	case '[':
		return Instruction{Op: OpLoopStart}, nil
	case ']':
		return Instruction{Op: OpLoopEnd}, nil
	}

	return Instruction{}, errors.New(fmt.Sprintf("Invalid instruction: %c", ch))
}
