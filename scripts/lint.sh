#!/bin/bash

# Lint script for ccexplorer
set -e

GOPATH=${GOPATH:-$(go env GOPATH)}
GOBIN=${GOBIN:-$GOPATH/bin}
GOLANGCI_LINT=${GOLANGCI_LINT:-$GOBIN/golangci-lint}

# Install golangci-lint if not present
if [ ! -f "$GOLANGCI_LINT" ]; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "Running linter..."
$GOLANGCI_LINT run --color=always --sort-results ./...

echo "Linting completed!"