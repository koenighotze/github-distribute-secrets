.DEFAULT_GOAL := build

.PHONY: all build test vet clean

install.tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest

clean:
	go clean -x -i

fmt:
	go fmt ./cmd/... ./internal/...

vet: fmt
	go vet ./cmd/... ./internal/...

lint:
	golangci-lint run ./...

deps.upgrade:
	go get -u ./...
	go mod tidy

deps.vendor:
	go mod vendor

deps.vulncheck:
	govulncheck ./...

deps.nancy:
	go list -json -deps ./... | docker run --rm -i sonatypecommunity/nancy:latest sleuth

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
