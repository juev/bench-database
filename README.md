# bench-database

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkInsert-8                      	       1	5479740375 ns/op
BenchmarkTransactionInsert-8           	       1	12250542875 ns/op
BenchmarkBulkInsert-8                  	      14	  77644524 ns/op
BenchmarkTransactionBulkInsert-8       	      14	  79995140 ns/op
BenchmarkCopyFromInsert-8              	      68	  19369895 ns/op
BenchmarkTransactionCopyFromInsert-8   	      64	  18720895 ns/op
BenchmarkUpdate-8                      	      43	  27245904 ns/op
BenchmarkUpdateWithTemporaryTable-8    	      34	  33986902 ns/op
PASS
ok  	github.com/juev/bench-database	28.028s
```
