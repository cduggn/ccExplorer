#!/bin/bash

# Test script for ccexplorer
set -e

echo "Running tests..."
go test ./...

echo "Running tests with race detection..."
go test -race ./...

echo "Running tests with coverage..."
go test -cover ./...

echo "All tests completed!"