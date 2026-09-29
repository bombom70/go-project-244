COVERAGE_MIN      = 70
COVERAGE_FILE     = coverage.out
COVERAGE_FILTERED = coverage.filtered.out

start:
	go run  cmd/gendiff/main.go

build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run

fix:
	golangci-lint run --fix

test:
	go test

test-coverage:
	go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@grep -vE 'main\.go:' $(COVERAGE_FILE) > $(COVERAGE_FILTERED) || cp $(COVERAGE_FILE) $(COVERAGE_FILTERED)
	@total=$$(go tool cover -func=$(COVERAGE_FILTERED) | awk '/^total:/ {gsub(/%/,"",$$3); print $$3}'); \
	echo "Total coverage: $$total% (min: $(COVERAGE_MIN)%)"; \
	if awk "BEGIN {exit !($$total < $(COVERAGE_MIN))}"; then \
		echo "❌ Coverage $$total% is below minimum $(COVERAGE_MIN)%"; \
		exit 1; \
	fi; \
	echo "✅ Coverage OK"
