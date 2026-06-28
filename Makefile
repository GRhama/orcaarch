.PHONY: test bench run-api export-csv

test:
	go test ./...

bench:
	go test ./internal/service/... -bench=. -benchmem -benchtime=3x

run-api:
	go run cmd/api/main.go

export-csv:
	go run cmd/export/main.go
