.PHONY: test bench run-api export-csv

VOLUME ?= 1000

test:
	go test ./...

bench:
	go test ./internal/service/... -bench=. -benchmem -benchtime=3x

run-api:
	go run cmd/api/main.go -volume $(VOLUME)

export-csv:
	go run cmd/export/main.go -volume $(VOLUME)
