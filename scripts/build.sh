#!/bin/bash

# Build script for ccexplorer
set -e

echo "Building ccexplorer..."
go env -w CGO_ENABLED=1
go build -o bin/ ./cmd/ccexplorer

echo "Build completed successfully!"