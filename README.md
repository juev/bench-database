# bench-database

```plain
goos: darwin
goarch: arm64
pkg: github.com/juev/bench-database
BenchmarkInsert-8                                 	       1	5054827458 ns/op
BenchmarkTransactionInsert-8                      	       1	11033566292 ns/op
BenchmarkBulkInsert-8                             	       9	 116801417 ns/op
BenchmarkTransactionBulkInsert-8                  	       9	 114834782 ns/op
BenchmarkCopyFromInsert-8                         	      31	  37822507 ns/op
BenchmarkTransactionCopyFromInsert-8              	      28	  38409595 ns/op
BenchmarkUpdate-8                                 	       1	77323176208 ns/op
BenchmarkUpdateWithTemporaryTable-8               	       1	15790238416 ns/op
BenchmarkUpdateWithTemporaryTableWithoutIndex-8   	       1	16610539250 ns/op
PASS
ok  	github.com/juev/bench-database	147.603s
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
