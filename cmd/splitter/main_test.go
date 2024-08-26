package main

import "testing"

func Benchmark_Main(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := run("../../data/huge.json", 2)
		if err != nil {
			b.Fatal(err)
		}
	}
}
