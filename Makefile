LOCAL_BIN:=$(CURDIR)/bin
GOOSE_BIN:=$(LOCAL_BIN)/goose

# test run benchmarks
.PHONY: test
test:
	@go test -bench=. -benchtime=10x