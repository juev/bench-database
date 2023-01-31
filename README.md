# bench-database

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkInsert-8                      	       1	4845428791 ns/op
BenchmarkTransactionInsert-8           	       1	11041694500 ns/op
BenchmarkBulkInsert-8                  	      14	  76598390 ns/op
BenchmarkTransactionBulkInsert-8       	      14	  77411000 ns/op
BenchmarkCopyFromInsert-8              	      67	  17556253 ns/op
BenchmarkTransactionCopyFromInsert-8   	      64	  18427501 ns/op
PASS
ok  	github.com/juev/bench-database	21.962s
```
