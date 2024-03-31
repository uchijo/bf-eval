package tool

import "github.com/uchijo/bf-eval/instr"

func Parse(input []uint8) ([]instr.Instruction, error) {
	retval := []instr.Instruction{}
	for _, c := range input {
		switch c {
		case '>', '<', '+', '-', '.', ',', '[', ']':
			inst, err := instr.NewInstruction(c)
			if err != nil {
				return nil, err
			}
			retval = append(retval, inst)
		}
	}
	return retval, nil
}
