
.DEFAULT_GOAL := build
.PHONY: build

GOBIN ?= $(shell go env GOPATH)/bin
GOIMPORTS ?= $(GOBIN)/goimports
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GORELEASER ?= $(GOBIN)/goreleaser

$(GOIMPORTS):
	go get golang.org/x/tools/cmd/goimports@latest

$(GORELEASER):
	go install github.com/goreleaser/goreleaser@latest

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: $(GOIMPORTS) $(GOLANGCI_LINT)

release: $(GORELEASER)
	$(GORELEASER) release --rm-dist