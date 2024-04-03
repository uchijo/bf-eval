package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/uchijo/bf-eval/instr"
	"github.com/uchijo/bf-eval/optimizer"
	"github.com/uchijo/bf-eval/tool"
)

type program struct {
	w          io.Writer
	localIndex int
	loopLabel  int
	labelStack []int
}

func NewProgram(w io.Writer) *program {
	return &program{
		w:          w,
		localIndex: 7,
	}
}

func main() {
	var content bytes.Buffer

	for _, f := range os.Args[1:] {
		b, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read file:", err)
			os.Exit(1)
		}
		content.Write(b)
	}

	if content.Len() == 0 {
		if _, err := io.Copy(&content, os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, "read stdin:", err)
			os.Exit(1)
		}
	}

	parsed, err := tool.Parse(content.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse:", err)
		os.Exit(1)
	}

	hoge := optimizer.SumShift(optimizer.SumIncrDecr(optimizer.ResetToZeroPattern(parsed)))

	buf := strings.Builder{}
	p := NewProgram(&buf)
	p.emitHeader()

	for _, i := range hoge {
		switch i.Op {
		case instr.OpShiftRight:
			p.emitMovePtr(i.Data, true)
		case instr.OpShiftLeft:
			p.emitMovePtr(i.Data, false)
		case instr.OpIncr:
			p.emitAdd(i.Data)
		case instr.OpDecr:
			p.emitSub(i.Data)
		case instr.OpOutput:
			p.emitDot()
		case instr.OpLoopStart:
			p.emitLoopStart()
		case instr.OpLoopEnd:
			p.emitLoopEnd()
		case instr.OpZeroReset:
			p.emitResetToZero()
		}
	}

	p.emitFooter()
	fmt.Println(buf.String())
}

func (p *program) emitHeader() {
	p.w.Write([]byte("define i32 @main() {\n"))
	p.w.Write([]byte("  %1 = alloca ptr, align 8\n"))
	p.w.Write([]byte("  %2 = alloca ptr, align 8\n"))
	p.w.Write([]byte("  %3 = call noalias ptr @calloc(i64 noundef 4096, i64 noundef 1) #3\n"))
	p.w.Write([]byte("  store ptr %3, ptr %1, align 8\n"))
	p.w.Write([]byte("  %4 = load ptr, ptr %1, align 8\n"))
	p.w.Write([]byte("  store ptr %4, ptr %2, align 8\n"))
	p.w.Write([]byte("  %5 = load ptr, ptr %2, align 8\n"))
	p.w.Write([]byte("  %6 = getelementptr inbounds i8, ptr %5, i64 2048\n"))
	p.w.Write([]byte("  store ptr %6, ptr %2, align 8\n"))
}

func (p *program) emitFooter() {
	p.w.Write([]byte("  ret i32 0\n}\n"))
	p.w.Write([]byte("declare noalias ptr @calloc(i64 noundef, i64 noundef) #1\n"))
	p.w.Write([]byte("declare i32 @putchar(i32 noundef) #2"))
	p.w.Write([]byte("attributes #1 = { nounwind allocsize(0,1) \"frame-pointer\"=\"all\" \"no-trapping-math\"=\"true\" \"stack-protector-buffer-size\"=\"8\" \"target-cpu\"=\"x86-64\" \"target-features\"=\"+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87\" \"tune-cpu\"=\"generic\" }\n"))
	p.w.Write([]byte("attributes #2 = { \"frame-pointer\"=\"all\" \"no-trapping-math\"=\"true\" \"stack-protector-buffer-size\"=\"8\" \"target-cpu\"=\"x86-64\" \"target-features\"=\"+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87\" \"tune-cpu\"=\"generic\" }"))
	p.localIndex += 3
}

func (p *program) emitMovePtr(offset int32, isAdd bool) {
	p.w.Write([]byte(fmt.Sprintf("%%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))
	if isAdd {
		p.w.Write([]byte(fmt.Sprintf("%%%v = getelementptr inbounds i8, ptr %%%v, i32 %v\n", p.localIndex+1, p.localIndex, offset)))
	} else {
		p.w.Write([]byte(fmt.Sprintf("%%%v = getelementptr inbounds i8, ptr %%%v, i32 -%v\n", p.localIndex+1, p.localIndex, offset)))
	}
	p.w.Write([]byte(fmt.Sprintf("store ptr %%%v, ptr %%2, align 8\n", p.localIndex+1)))
	p.localIndex += 2
}

func (p *program) emitAdd(offset int32) {
	// load
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))

	// add
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 4\n", p.localIndex+1, p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = add i8 %%%v, %v\n", p.localIndex+2, p.localIndex+1, offset)))

	// save
	p.w.Write([]byte(fmt.Sprintf("  store i8 %%%v, ptr %%%v, align 1\n\n", p.localIndex+2, p.localIndex)))

	p.localIndex += 3
}

func (p *program) emitSub(offset int32) {
	// load
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))

	// add
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 4\n", p.localIndex+1, p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = sub i8 %%%v, %v\n", p.localIndex+2, p.localIndex+1, offset)))

	// save
	p.w.Write([]byte(fmt.Sprintf("  store i8 %%%v, ptr %%%v, align 1\n\n", p.localIndex+2, p.localIndex)))

	p.localIndex += 3
}

func (p *program) emitDot() {
	p.w.Write([]byte(";dot\n"))
	p.w.Write([]byte(fmt.Sprintf("%%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("%%%v = load i8, ptr %%%v, align 1\n", p.localIndex+1, p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("%%%v = sext i8 %%%v to i32\n", p.localIndex+2, p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("%%%v = call i32 @putchar(i32 noundef %%%v)\n", p.localIndex+3, p.localIndex+2)))
	p.localIndex += 4
}

func (p *program) emitLoopStart() {
	p.labelStack = append(p.labelStack, p.loopLabel)
	p.w.Write([]byte(fmt.Sprintf("  br label %%loop_start_%v\n", p.loopLabel)))

	// 最初の一回をやるか判断
	p.w.Write([]byte(fmt.Sprintf("\nloop_start_%v:\n", p.loopLabel)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 1\n", p.localIndex+1, p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = icmp ne i8 %%%v, 0\n", p.localIndex+2, p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  br i1 %%%v, label %%loop_body_%v, label %%loop_end_%v\n", p.localIndex+2, p.loopLabel, p.loopLabel)))

	p.w.Write([]byte(fmt.Sprintf("\nloop_body_%v:\n\n", p.loopLabel)))
	p.loopLabel++
	p.localIndex += 3
}

func (p *program) emitLoopEnd() {
	label := p.labelStack[len(p.labelStack)-1]
	p.labelStack = p.labelStack[:len(p.labelStack)-1]

	// 最初に戻るか判断
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 1\n", p.localIndex+1, p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = icmp eq i8 %%%v, 0\n", p.localIndex+2, p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  br i1 %%%v, label %%loop_end_%v, label %%loop_body_%v\n", p.localIndex+2, label, label)))
	p.w.Write([]byte(fmt.Sprintf("\nloop_end_%v:\n\n", label)))
	p.localIndex += 3
}

func (p *program) emitResetToZero() {
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  store i8 0, ptr %%%v, align 1\n", p.localIndex)))
	p.localIndex++
}
