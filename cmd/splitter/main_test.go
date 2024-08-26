package main

import "testing"

func Benchmark_Main(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run("../true.txt", 1)
	}
}
