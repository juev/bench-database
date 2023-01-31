# bench-database

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkInsert-8                      	       1	5342260083 ns/op
BenchmarkTransactionInsert-8           	       1	11400897500 ns/op
BenchmarkBulkInsert-8                  	      14	  75702506 ns/op
BenchmarkTransactionBulkInsert-8       	      14	  78499673 ns/op
BenchmarkCopyFromInsert-8              	      61	  21757240 ns/op
BenchmarkTransactionCopyFromInsert-8   	      44	  24599059 ns/op
BenchmarkUpdate-8                      	       1	29261891625 ns/op
BenchmarkUpdateWithTemporaryTable-8    	      39	  32260781 ns/op
PASS
ok  	github.com/juev/bench-database	54.126s
```
