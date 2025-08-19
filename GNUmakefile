# SPECIAL: Hardcoded provider prefix required in this file
REPO_BUILD_TAG_VAL = repoAWX

default: fmt lint install generate

build:
	go build -v -tags=$(REPO_BUILD_TAG_VAL) ./...

install: build
	go install -v -tags=$(REPO_BUILD_TAG_VAL) ./...

lint:
	golangci-lint run --build-tags=$(REPO_BUILD_TAG_VAL)

generate:
	GOFLAGS="-tags=$(REPO_BUILD_TAG_VAL)" go run generate-examples/main.go
	(cd tools && GOFLAGS="-tags=$(REPO_BUILD_TAG_VAL)" go generate ./...)

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 -tags=$(REPO_BUILD_TAG_VAL) ./internal/provider

testacc:
	TF_ACC=1 go test -tags=$(REPO_BUILD_TAG_VAL) -v -cover ./internal/provider

.PHONY: fmt lint test testacc build install generate
