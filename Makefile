# This file contains convenience targets for the project.
# It is not intended to be used as a build system.

.DEFAULT_GOAL := build

GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
GORELEASER ?= $(GOBIN)/goreleaser
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint

$(GORELEASER):
	go install github.com/goreleaser/goreleaser@latest

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: setup
setup:
	git config --local core.hooksPath .githooks/

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --color=always --sort-results ./...

.PHONY: lint-exp
lint-exp: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix --config .golangci-exp.yaml ./...

.PHONY: test-race
test-race:
	go run test -race ./...

.PHONY: test-cover
test-cover:
	go run test -cover ./...

.PHONY: run-pkgsite
run-pkgsite:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest

goimports: $(GOIMPORTS)
	$(GOIMPORTS) -l ./cmd ./internal

.PHONY: run
run:
	go env -w CGO_ENABLED=1
	go run ./cmd/cloudcost/cloudcost.go

.PHONY: build
build:
	go build -o bin/ ./...

.PHONY: release
release: $(GORELEASER)
	$(GORELEASER) release --rm-dist

.PHONY: clean
clean: clean-lint-cache

.PHONY: clean-lint-cache
clean-lint-cache:
	golangci-lint cache clean

.PHONY: git-secrets
git-secrets:
	git secrets --register-aws --global

.PHONY: tag   # make tag VERSION=0.6.0
tag:
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	#git tag -a v0.6.0 <commit-id>  -m "Release v0.6.0"
	git push origin v$(VERSION)
