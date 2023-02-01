# bench-database

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkInsert-8                                 	       1	5134397000 ns/op
BenchmarkTransactionInsert-8                      	       1	10928278792 ns/op
BenchmarkBulkInsert-8                             	       9	 114961366 ns/op
BenchmarkTransactionBulkInsert-8                  	       9	 115041338 ns/op
BenchmarkCopyFromInsert-8                         	      31	  37682081 ns/op
BenchmarkTransactionCopyFromInsert-8              	      30	  38555136 ns/op
BenchmarkBulkUpdate-8                             	       1	64325916625 ns/op
BenchmarkUpdateWithTemporaryTable-8               	       1	15860838709 ns/op
BenchmarkUpdateWithTemporaryTableWithoutIndex-8   	       1	15335684083 ns/op
PASS
ok  	github.com/juev/bench-database	133.811s
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
