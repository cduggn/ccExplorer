
.DEFAULT_GOAL := build
.PHONY: build checks imports release lint setup

GO ?= go
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOIMPORTS ?= $(GOBIN)/goimports
STATICCHECK ?= $(GOBIN)/staticcheck
GORELEASER ?= $(GOBIN)/goreleaser
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint

setup:
	git config --local core.hooksPath .githooks/

$(GOIMPORTS):
	$(GO) install golang.org/x/tools/cmd/goimports@latest

$(GORELEASER):
	$(GO) install github.com/goreleaser/goreleaser@latest

$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

$(GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

imports: $(GOIMPORTS)
	$(GOIMPORTS) -l cmd/ internal/ && echo "OK"

checks: $(STATICCHECK) lint imports
	$(GO) vet ./...
	$(STATICCHECK) ./...
	$(GOLANGCI_LINT) run ./...

build: checks
	$(GO) build -o bin/ ./...

release: $(GORELEASER)
	$(GORELEASER) release --rm-dist