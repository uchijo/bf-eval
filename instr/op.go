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
	Op Op
}

func NewInstruction(ch uint8) (Instruction, error) {
	switch ch {
	case '>':
		return Instruction{OpShiftRight}, nil
	case '<':
		return Instruction{OpShiftLeft}, nil
	case '+':
		return Instruction{OpIncr}, nil
	case '-':
		return Instruction{OpDecr}, nil
	case '.':
		return Instruction{OpOutput}, nil
	case ',':
		return Instruction{OpInput}, nil
	case '[':
		return Instruction{OpLoopStart}, nil
	case ']':
		return Instruction{OpLoopEnd}, nil
	}

	return Instruction{}, errors.New(fmt.Sprintf("Invalid instruction: %c", ch))
}
