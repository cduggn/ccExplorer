# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

**Build and Test:**
- `make build` - Build application to bin/ directory
- `make test` - Run all tests
- `make test-race` - Run tests with race detection
- `make test-cover` - Run tests with coverage
- `make lint` - Run golangci-lint
- `make run` - Run application directly

**Release:**
- `make release` - Create release using goreleaser

## Architecture Overview

ccExplorer is a CLI tool for AWS cost analysis that follows clean architecture patterns with clear separation of concerns.

**Core Data Flow:**
1. CLI layer (`cmd/cli/`) parses commands and flags using Cobra
2. Commands synthesize requests to internal domain types (`internal/types/`)
3. AWS service layer (`internal/awsservice/`) calls Cost Explorer API
4. Utils (`internal/utils/`) transform AWS responses to internal types
5. Writer layer (`internal/writer/`) formats output (stdout, CSV, charts, Pinecone vector DB)

**Key Architectural Patterns:**
- **Dependency Inversion**: Interfaces defined in `internal/ports/` 
- **Factory Pattern**: Writer creation based on output format in `internal/writer/`
- **Strategy Pattern**: Multiple sorting and filtering implementations
- **Adapter Pattern**: AWS SDK integration through service abstraction

## Package Responsibilities

**`internal/awsservice/`** - AWS Cost Explorer client wrapper implementing `ports.AWSService`
**`internal/writer/`** - Output formatting with multiple `Printer` implementations (stdout, CSV, charts, Pinecone)
**`internal/types/`** - Domain models for requests, responses, and presentation
**`internal/utils/`** - Data transformations, sorting, date handling
**`internal/openai/`** - Embedding generation for vector database storage
**`internal/pinecone/`** - Vector database client with batch processing
**`internal/ports/`** - Interface definitions for testability and modularity

## Key Integration Points

**AWS Integration:** Uses AWS SDK v2 with Cost Explorer API. Authentication handled through standard AWS credential chain.

**AI/Vector DB Pipeline:** Cost data → OpenAI embeddings → Pinecone storage for semantic search capabilities.

**CLI Structure:** Main commands are `ccexplorer get aws` (cost queries) and `ccexplorer get aws forecast` (forecasting). Supports complex filtering by AWS dimensions and cost allocation tags.

## Testing Notes

Tests are located alongside source files. The codebase includes unit tests for flag parsing (`cmd/cli/flags/*_test.go`), utility functions (`internal/utils/commons_test.go`), and service integrations (`internal/awsservice/filter_test.go`, `internal/pinecone/pinecone_test.go`).

# Additional Instructions
gemini.md
