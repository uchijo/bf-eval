package tool

import (
	"fmt"

	"github.com/uchijo/bf-eval/instr"
)

func DumpInstr(src []instr.Instruction) {
	nest := 0
	for _, i := range src {
		switch i.Op {
		case instr.OpShiftRight:
			fmt.Printf("%vShiftRight %v\n", padding(nest), i.Data)
		case instr.OpShiftLeft:
			fmt.Printf("%vShiftLeft %v\n", padding(nest), i.Data)
		case instr.OpIncr:
			fmt.Printf("%vIncr %v\n", padding(nest), i.Data)
		case instr.OpDecr:
			fmt.Printf("%vDecr %v\n", padding(nest), i.Data)
		case instr.OpOutput:
			fmt.Printf("%vOutput\n", padding(nest))
		case instr.OpInput:
			fmt.Printf("%vInput\n", padding(nest))
		case instr.OpZeroReset:
			fmt.Printf("%vZeroReset\n", padding(nest))
		case instr.OpLoopStart:
			fmt.Printf("%vLoopStart\n", padding(nest))
			nest++
		case instr.OpLoopEnd:
			nest--
			fmt.Printf("%vLoopEnd\n", padding(nest))
		case instr.OpAddMem:
			fmt.Printf("%vAddMem %v\n", padding(nest), i.Data)
		case instr.OpMultiShift:
			fmt.Printf("%vMultiShift %v\n", padding(nest), i.Data)
		}
	}
}

func padding(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "  "
	}
	return s
}
