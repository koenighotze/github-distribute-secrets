.DEFAULT_GOAL := build

.PHONY: all build test vet clean

clean:
	go clean -x -i

fmt:
	go fmt ./cmd/... ./internal/...

vet: fmt
	go vet ./cmd/... ./internal/...

get.dependencies:
	go mod tidy

test: get.dependencies
	echo TODO

test.report: get.dependencies
	go test -json > TestResults.json

build: vet get.dependencies
	go build -o github-distribute-secrets ./cmd/github-distribute-secrets

run.local:
	go run ./cmd/github-distribute-secrets/main.go
