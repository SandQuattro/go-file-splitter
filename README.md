# Concurrent chunks file splitter

Using concurrency and sync.Pool

Benching 26Mb json file

go test -v -bench=. -benchmem

14 Core, 36Gb Mem
```
2024/08/26 23:08:18 file ../../data/huge.json size 26141343
2024/08/26 23:08:18 chunk 1  read 2097152 bytes
2024/08/26 23:08:18 chunk 1 write 2097152 bytes
...
2024/08/26 23:08:19 chunk 13  read 975519 bytes
2024/08/26 23:08:19 chunk 13 write 975519 bytes
Benchmark_Main-14   184           6075444 ns/op         2161932 B/op        171 allocs/op
PASS
ok      file_splitter/cmd/splitter      3.111s
```