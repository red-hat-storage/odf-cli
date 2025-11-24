build:
	gofmt -w $(shell find . -type f -name '*.go')
	@echo
	env GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o bin/odf cmd/odf/main.go

clean:
	@rm -f bin/odf

test:
	@echo "running unit tests"
	go test ./...

help:
	@echo "build : Create go binary."
	@echo "test  : Runs unit tests"
	@echo "clean : Remove go binary file."
