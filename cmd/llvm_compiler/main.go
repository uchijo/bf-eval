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

	parseed, err := tool.Parse(content.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse:", err)
		os.Exit(1)
	}
	optimized := optimizer.Optimize(parseed)

	buf := strings.Builder{}
	p := NewProgram(&buf)
	p.emitHeader()

	for _, i := range optimized {
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
		}
	}

	p.emitFooter()
	fmt.Println(buf.String())
}

func (p *program) emitHeader() {
	p.w.Write([]byte("@stdout = external global ptr, align 8\n"))
	p.w.Write([]byte("define i32 @main() {\n"))
	p.w.Write([]byte("  %1 = alloca ptr, align 8\n"))
	p.w.Write([]byte("  %2 = alloca ptr, align 8\n"))
	p.w.Write([]byte("  %3 = alloca i32, align 4\n"))
	p.w.Write([]byte("  %4 = alloca i32, align 4\n"))
	p.w.Write([]byte("  %5 = call noalias ptr @calloc(i64 noundef 4096, i64 noundef 1) #3\n"))
	p.w.Write([]byte("  store ptr %5, ptr %1, align 8\n"))
	p.w.Write([]byte("  %6 = call noalias ptr @calloc(i64 noundef 8192, i64 noundef 1) #3\n"))
	p.w.Write([]byte("  store ptr %6, ptr %2, align 8\n"))
	p.w.Write([]byte("  store i32 2048, ptr %3, align 4\n"))
	p.w.Write([]byte("  store i32 0, ptr %4, align 4\n"))
}

func (p *program) emitFooter() {
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr @stdout, align 8\n", p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = call i32 @fputs(ptr noundef %%%v, ptr noundef %%%v\n)", p.localIndex+2, p.localIndex, p.localIndex+1)))
	p.w.Write([]byte("  ret i32 0\n}\n"))
	p.w.Write([]byte("declare noalias ptr @calloc(i64 noundef, i64 noundef) #1\n"))
	p.w.Write([]byte("declare i32 @fputs(ptr noundef, ptr noundef) #2\n"))
	p.w.Write([]byte("attributes #1 = { nounwind allocsize(0,1) \"frame-pointer\"=\"all\" \"no-trapping-math\"=\"true\" \"stack-protector-buffer-size\"=\"8\" \"target-cpu\"=\"x86-64\" \"target-features\"=\"+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87\" \"tune-cpu\"=\"generic\" }\n"))
	p.w.Write([]byte("attributes #2 = { \"frame-pointer\"=\"all\" \"no-trapping-math\"=\"true\" \"stack-protector-buffer-size\"=\"8\" \"target-cpu\"=\"x86-64\" \"target-features\"=\"+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87\" \"tune-cpu\"=\"generic\" }"))
	p.localIndex += 3
}

func (p *program) emitMovePtr(offset int32, isAdd bool) {
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i32, ptr %%3, align 4\n", p.localIndex)))
	if isAdd {
		p.w.Write([]byte(fmt.Sprintf("  %%%v = add i32 %%%v, %v\n", p.localIndex+1, p.localIndex, offset)))
	} else {
		p.w.Write([]byte(fmt.Sprintf("  %%%v = sub i32 %%%v, %v\n", p.localIndex+1, p.localIndex, offset)))
	}
	p.w.Write([]byte(fmt.Sprintf("  store i32 %%%v, ptr %%3, align 4\n", p.localIndex+1)))
	p.localIndex += 2
}

func (p *program) emitAdd(offset int32) {
	// load
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%1, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i32, ptr %%3, align 4\n", p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = zext i32 %%%v to i64\n", p.localIndex+2, p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = getelementptr inbounds i8, ptr %%%v, i64 %%%v\n", p.localIndex+3, p.localIndex, p.localIndex+2)))

	// add
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 1\n", p.localIndex+4, p.localIndex+3)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = add i8 %%%v, %v\n", p.localIndex+5, p.localIndex+4, offset)))

	// save
	p.w.Write([]byte(fmt.Sprintf("  store i8 %%%v, ptr %%%v, align 1\n", p.localIndex+5, p.localIndex+3)))

	p.localIndex += 6
}

func (p *program) emitSub(offset int32) {
	// load
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%1, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i32, ptr %%3, align 4\n", p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = zext i32 %%%v to i64\n", p.localIndex+2, p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = getelementptr inbounds i8, ptr %%%v, i64 %%%v\n", p.localIndex+3, p.localIndex, p.localIndex+2)))

	// sub
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 1\n", p.localIndex+4, p.localIndex+3)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = sub i8 %%%v, %v\n", p.localIndex+5, p.localIndex+4, offset)))

	// save
	p.w.Write([]byte(fmt.Sprintf("  store i8 %%%v, ptr %%%v, align 1\n", p.localIndex+5, p.localIndex+3)))

	p.localIndex += 6
}

func (p *program) emitDot() {
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%1, align 8\n", p.localIndex)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i32, ptr %%3, align 4\n", p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = zext i32 %%%v to i64\n", p.localIndex+2, p.localIndex+1)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = getelementptr inbounds i8, ptr %%%v, i64 %%%v\n", p.localIndex+3, p.localIndex, p.localIndex+2)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i8, ptr %%%v, align 1\n", p.localIndex+4, p.localIndex+3)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load ptr, ptr %%2, align 8\n", p.localIndex+5)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i32, ptr %%4, align 4\n", p.localIndex+6)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = zext i32 %%%v to i64\n", p.localIndex+7, p.localIndex+6)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = getelementptr inbounds i8, ptr %%%v, i64 %%%v\n", p.localIndex+8, p.localIndex+5, p.localIndex+7)))
	p.w.Write([]byte(fmt.Sprintf("  store i8 %%%v, ptr %%%v, align 1\n", p.localIndex+4, p.localIndex+8)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = load i32, ptr %%4, align 4\n", p.localIndex+9)))
	p.w.Write([]byte(fmt.Sprintf("  %%%v = add i32 %%%v, 1\n", p.localIndex+10, p.localIndex+9)))
	p.w.Write([]byte(fmt.Sprintf("  store i32 %%%v, ptr %%4, align 4\n", p.localIndex+10)))
	p.localIndex += 11
}
