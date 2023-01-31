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

## update table

We have limit records on batch, when we use simple update. In this case 3_000.

But when we use temporary tables we don`t have this limits.

With 3_000 elements:

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkUpdateWithTemporaryTable-8   	      34	  34589695 ns/op
PASS
ok  	github.com/juev/bench-database	2.145s
```

With 30_000_000 elements:

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkUpdateWithTemporaryTable-8   	       1	46225218167 ns/op
PASS
ok  	github.com/juev/bench-database	47.001s
```
