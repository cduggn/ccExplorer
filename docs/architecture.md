# Architecture Documentation

## Overview

ccExplorer follows clean architecture patterns with clear separation of concerns.

## Directory Structure

```
ccExplorer/
├── cmd/ccexplorer/     # Main application entry point
├── internal/           # Private application code
│   ├── awsservice/     # AWS Cost Explorer integration
│   ├── writer/         # Output formatting
│   ├── types/          # Domain models
│   ├── utils/          # Utilities and transformations
│   ├── ports/          # Interface definitions
│   └── ...
├── configs/            # Configuration examples
├── scripts/            # Build and utility scripts
├── docs/               # Documentation
├── api/                # API definitions
├── examples/           # Usage examples
└── tools/              # Supporting tools
```

## Data Flow

1. CLI layer parses commands using Cobra
2. Commands create domain requests
3. AWS service layer calls Cost Explorer API
4. Utils transform AWS responses to internal types
5. Writer layer formats output

## Key Patterns

- **Dependency Inversion**: Interfaces in `/internal/ports/`
- **Factory Pattern**: Writer creation in `/internal/writer/`
- **Strategy Pattern**: Multiple sorting/filtering implementations
- **Adapter Pattern**: AWS SDK integration