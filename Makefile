.DEFAULT_GOAL := build

.PHONY: all build test vet clean install.tools lint deps.vulncheck

install.tools: install.tools.local

install.tools.local:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	brew install golangci-lint

clean:
	go clean -x -i
	rm -f ./github-distribute-secrets

fmt:
	go fmt ./cmd/... ./internal/... ./pkg/...

vet: fmt
	go vet ./cmd/... ./internal/... ./pkg/...

lint: lint.local

lint.local:
	golangci-lint run ./...

deps.upgrade:
	go get -u ./...
	go mod tidy

deps.vendor:
	go mod vendor

deps.vulncheck: deps.vulncheck.local

deps.vulncheck.local:
	govulncheck ./...

deps.nancy:
	go list -json -deps ./... | docker run --rm -i sonatypecommunity/nancy:latest sleuth

get.dependencies:
	go mod tidy

test: get.dependencies
	go test ./internal/... ./cmd/... ./pkg/... -coverprofile=coverage.out

test.all: get.dependencies
	go test -tags=integration ./internal/... ./cmd/... ./pkg/... -coverprofile=coverage.out

test.report: test
	go tool cover -html=coverage.out

build: get.dependencies
	go build -o github-distribute-secrets ./cmd/github-distribute-secrets

run.local:
	go run ./cmd/github-distribute-secrets/main.go
