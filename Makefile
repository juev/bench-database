.PHONY: start-db
start-db:
	@docker compose -f docker/docker-compose.yml up -d

.PHONY: stop-db
stop-db:
	@docker compose -f docker/docker-compose.yml down

# test run benchmarks
.PHONY: test
test:
	@go test -bench=.
