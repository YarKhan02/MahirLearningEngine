.PHONY: lint test build security run

lint:
	go vet ./...
	golangci-lint run

test:
	go test ./... -race

build:
	go build ./...

security:
	govulncheck ./...

run:
	go run ./cmd/server
