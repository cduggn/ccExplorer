
.DEFAULT_GOAL := build
.PHONY: build checks imports release lint setup

GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOIMPORTS ?= $(GOBIN)/goimports
STATICCHECK ?= $(GOBIN)/staticcheck
GORELEASER ?= $(GOBIN)/goreleaser
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint

setup:
	git config --local core.hooksPath .githooks/

$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@latest

$(GORELEASER):
	go install github.com/goreleaser/goreleaser@latest

$(STATICCHECK):
	go install honnef.co/go/tools/cmd/staticcheck@latest

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

imports: $(GOIMPORTS)
#	$(GOIMPORTS) -l cmd/ internal/ && echo "OK"

checks: $(STATICCHECK) lint #imports
	go vet ./...
	$(STATICCHECK) ./...
	$(GOLANGCI_LINT) run ./...

build: checks
	go build -o bin/ ./...

release: $(GORELEASER)
	$(GORELEASER) release --rm-dist