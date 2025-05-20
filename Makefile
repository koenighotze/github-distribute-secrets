.DEFAULT_GOAL := build

.PHONY: all build test vet clean

clean:
	go clean -x -i

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

get.dependencies:
	go get .

test: get.dependencies
	echo TODO

test.report: get.dependencies
	go test -json > TestResults.json

build: vet get.dependencies
	go build -o cmd/

run.local:
	go run gh-distribute-secrets.go
