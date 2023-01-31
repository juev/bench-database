# bench-database

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkInsert-8                                 	       1	5030175458 ns/op
BenchmarkTransactionInsert-8                      	       1	11084082667 ns/op
BenchmarkBulkInsert-8                             	      14	  76770756 ns/op
BenchmarkTransactionBulkInsert-8                  	      15	  79860147 ns/op
BenchmarkCopyFromInsert-8                         	      62	  17539187 ns/op
BenchmarkTransactionCopyFromInsert-8              	      58	  18794177 ns/op
BenchmarkUpdate-8                                 	       1	6247260000 ns/op
BenchmarkUpdateWithTemporaryTable-8               	       1	3884489625 ns/op
BenchmarkUpdateWithTemporaryTableWithoutIndex-8   	       1	3289993084 ns/op
PASS
ok  	github.com/juev/bench-database	35.535s
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
