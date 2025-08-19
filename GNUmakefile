# SPECIAL: Hardcoded provider prefix required in this file
default: fmt lint install generate

build:
	go build -v -tags=repoAAP ./...

install: build
	go install -v -tags=repoAAP ./...

lint:
	golangci-lint run --build-tags=repoAAP

generate:
	GOFLAGS="-tags=repoAAP" go run generate-examples/main.go
	(cd tools && GOFLAGS="-tags=repoAAP" go generate ./...)

fmt:
	gofmt -s -w -e .

test:
# SPECIAL: add -tags=repoAWX or -tags=repoAAP
	go test -v -cover -timeout=120s -parallel=10 -tags=repoAAP ./internal/provider

testacc:
# SPECIAL: add -tags=repoAWX or -tags=repoAAP
	TF_ACC=1 go test -tags=repoAAP -v -cover ./internal/provider

.PHONY: fmt lint test testacc build install generate
