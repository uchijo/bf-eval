optimized-bf: optimized-bf.s
	gcc -o optimized-bf optimized-bf.s

optimized-bf.s: optimized-bf.ll
	llc-18 -O3 -o optimized-bf.s optimized-bf.ll

optimized-bf.ll: compiled-bf.ll
	opt-18 -O3 -o optimized-bf.ll compiled-bf.ll

compiled-bf.ll: m.bf
	go run cmd/llvm_compiler/main.go m.bf > compiled-bf.ll

clean:
	rm -f optimized-bf optimized-bf.s optimized-bf.ll compiled-bf.ll

.PHONY: clean
