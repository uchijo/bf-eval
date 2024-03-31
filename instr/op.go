package instr

import (
	"errors"
	"fmt"
)

type Op uint8

const (
	OpShiftRight Op = iota
	OpShiftLeft
	OpIncr
	OpDecr
	OpOutput
	OpInput
	OpZeroReset
	OpLoopStart
	OpLoopEnd
)

type Instruction struct {
	Op   Op
	Data uint8
}

func NewInstruction(ch uint8) (Instruction, error) {
	switch ch {
	case '>':
		return Instruction{Op: OpShiftRight, Data: 1}, nil
	case '<':
		return Instruction{Op: OpShiftLeft, Data: 1}, nil
	case '+':
		return Instruction{Op: OpIncr}, nil
	case '-':
		return Instruction{Op: OpDecr}, nil
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
