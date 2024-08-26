package main

import "testing"

func Benchmark_Main(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := run("../true.txt", 1)
		if err != nil {
			b.Fatal(err)
		}
	}
}
