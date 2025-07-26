# Add GCP Cloud Billing API integration for multi-cloud cost analysis

**Labels:** `enhancement`, `gcp`, `multi-cloud`, `feature`  
**Priority:** Medium  
**Complexity:** High  

## Problem Statement

ccExplorer currently provides comprehensive AWS cost analysis through the Cost Explorer API, but lacks support for Google Cloud Platform cost data. Organizations using multi-cloud architectures need unified cost visibility across AWS and GCP to make informed financial decisions and optimize cloud spending holistically.

The current architecture is well-positioned for multi-cloud expansion due to its clean architecture patterns, interface-driven design, and generic utilities, but requires extension to support GCP's Cloud Billing API ecosystem.

## Acceptance Criteria

### Core Functionality
- [ ] **GCP Service Interface**: Define `GCPService` interface in `internal/ports/services.go` following existing `AWSService` patterns
- [ ] **GCP Service Implementation**: Create `internal/gcpservice/` package with client implementing the interface
- [ ] **CLI Commands**: Add `ccexplorer get gcp` command structure parallel to existing `ccexplorer get aws`
- [ ] **Multi-format Output**: Support all existing output formats (stdout, CSV, charts, Pinecone) for GCP data
- [ ] **Date Range Support**: Implement equivalent date filtering and granularity options
- [ ] **Resource Filtering**: Support GCP-specific dimensions (project, service, SKU, region, etc.)
- [ ] **Cost Metrics**: Support GCP billing data types (cost, usage, credits, discounts)

### Authentication & Configuration
- [ ] **Service Account Support**: Support GCP service account authentication
- [ ] **Application Default Credentials**: Fall back to ADC when service account not provided
- [ ] **Configuration Management**: Extend existing config system for GCP credentials
- [ ] **Environment Variables**: Support `GOOGLE_APPLICATION_CREDENTIALS` and project ID configuration

### Data Integration
- [ ] **Type Mapping**: Map GCP billing data to existing internal types where possible
- [ ] **GCP-specific Types**: Extend type system for GCP-unique concepts (projects, billing accounts)
- [ ] **Data Transformation**: Transform GCP API responses to internal representation
- [ ] **Error Handling**: Implement GCP-specific error handling and API limits

### Testing & Quality
- [ ] **Unit Tests**: Comprehensive unit tests for GCP service layer
- [ ] **Integration Tests**: Tests against GCP Billing API (with proper mocking)
- [ ] **Flag Parsing Tests**: Tests for GCP-specific command flags
- [ ] **Data Transformation Tests**: Validate GCP data mapping accuracy

## Technical Approach

### 1. Interface Extension
Extend `internal/ports/services.go` to include GCP service interface:

```go
type GCPService interface {
    GetBillingData(ctx context.Context, req types.GCPBillingRequest) (*billing.ListServicesResponse, error)
    GetCostData(ctx context.Context, req types.GCPCostRequest) (*billing.ProjectBillingInfo, error)
    GetUsageData(ctx context.Context, req types.GCPUsageRequest) (*cloudasset.SearchAllResourcesResponse, error)
}
```

### 2. Service Implementation Architecture
Create `internal/gcpservice/` package following AWS patterns:

```
internal/gcpservice/
├── client.go              # Main service implementation
├── config.go              # GCP authentication configuration  
├── get_billing_data.go     # Billing API integration
├── get_cost_data.go        # Cost data retrieval
├── get_usage_data.go       # Usage data retrieval
├── filters.go              # GCP-specific filtering logic
└── client_test.go          # Service layer tests
```

### 3. CLI Integration
Extend `cmd/cli/get_command.go` with GCP command structure:

```go
type GCPCommandType struct {
    Cmd *cobra.Command
}

func (g *GCPCommandType) DefineFlags() {
    // Project ID (required)
    g.Cmd.Flags().StringVarP(&gcpProjectID, "project", "p", "", "GCP Project ID (required)")
    g.Cmd.MarkFlagRequired("project")
    
    // Billing account filter
    g.Cmd.Flags().StringVarP(&gcpBillingAccount, "billingAccount", "b", "", "Filter by billing account")
    
    // GCP-specific dimensions (service, SKU, project, region, etc.)
    // Date ranges, granularity, output format (reuse existing patterns)
}
```

### 4. Data Type Extensions
Extend `internal/types/command_types.go` with GCP-specific types:

```go
type GCPBillingRequest struct {
    ProjectID       string
    BillingAccount  string
    Time           Time
    Granularity    string
    Services       []string
    Regions        []string
    PrintFormat    string
    // ... other fields following AWS patterns
}

type GCPBillingResponse struct {
    BillingData    map[string]GCPBillingItem
    ProjectInfo    string
    TotalCost      float64
    Currency       string
    // ... following internal type patterns
}
```

## Security Considerations

### Authentication Security
- **Service Account Keys**: Store service account keys securely, never in version control
- **IAM Principles**: Follow principle of least privilege for GCP service accounts
- **Required Permissions**: Document minimum IAM roles needed (`billing.viewer`, `cloudasset.viewer`)
- **Credential Rotation**: Support credential rotation through standard GCP mechanisms

### API Security
- **Rate Limiting**: Implement proper rate limiting for GCP Billing API calls
- **Request Validation**: Validate all user inputs for project IDs and billing accounts
- **Error Sanitization**: Sanitize GCP API errors to avoid information disclosure
- **Audit Logging**: Log GCP API access attempts for security monitoring

### Data Protection
- **Data Encryption**: Ensure billing data is encrypted in transit and at rest
- **PII Handling**: Handle potential PII in resource names/labels appropriately
- **Access Controls**: Verify user has appropriate GCP billing permissions

## Dependencies

### External Libraries
```go
// Primary GCP SDKs
cloud.google.com/go/billing/apiv1     v1.18.0  // Cloud Billing API
cloud.google.com/go/asset/apiv1       v1.19.0  // Cloud Asset Inventory API
google.golang.org/api/cloudbilling/v1 v0.0.0   // Additional billing capabilities

// Authentication and core
google.golang.org/api/option          v0.0.0   // Client options
golang.org/x/oauth2/google           v0.0.0   // OAuth2 for GCP
google.golang.org/grpc               v1.60.0  // gRPC support
```

### Integration Dependencies
- Update `go.mod` with GCP SDK dependencies
- Extend existing `internal/config/config.go` for GCP configuration
- Leverage existing `internal/utils/` transformation utilities
- Reuse `internal/writer/` output formatting system

## Implementation Plan

### Phase 1: Foundation (Week 1-2)
1. **Dependencies Setup**
   - Add GCP SDK dependencies to `go.mod`
   - Update build configurations and CI/CD

2. **Interface Design**
   - Define `GCPService` interface in `internal/ports/services.go`
   - Extend type system with GCP-specific types
   - Design authentication configuration structure

3. **Basic Client Implementation**
   - Implement basic GCP client in `internal/gcpservice/client.go`
   - Add authentication configuration handling
   - Create initial service registration patterns

### Phase 2: Core Integration (Week 3-4)
1. **Billing API Integration**
   - Implement `GetBillingData` method using Cloud Billing API
   - Add support for project-level cost queries
   - Implement basic filtering and date range support

2. **CLI Command Structure**
   - Add `ccexplorer get gcp` command with basic flags
   - Implement input validation and flag parsing
   - Create request synthesis following AWS patterns

3. **Data Transformation**
   - Map GCP billing responses to internal types
   - Implement cost aggregation and grouping logic
   - Add currency and unit conversions

### Phase 3: Feature Parity (Week 5-6)
1. **Advanced Filtering**
   - Support GCP dimensions (service, SKU, region, labels)
   - Implement complex filtering expressions
   - Add support for multiple projects

2. **Output Format Support**
   - Ensure stdout output works with GCP data
   - Test CSV export functionality
   - Validate chart generation with GCP cost data
   - Test Pinecone integration for GCP vectors

3. **Error Handling & Resilience**
   - Implement comprehensive error handling
   - Add retry logic for transient API failures
   - Create meaningful error messages for users

### Phase 4: Testing & Documentation (Week 7-8)
1. **Comprehensive Testing**
   - Unit tests for all GCP service methods
   - Integration tests with mock GCP APIs
   - CLI flag parsing and validation tests
   - End-to-end workflow testing

2. **Documentation**
   - Update README with GCP setup instructions
   - Document required GCP permissions
   - Add GCP command usage examples
   - Update architecture documentation

3. **Performance & Security Review**
   - Benchmark GCP API performance
   - Security review of authentication handling
   - Code review focusing on error handling

## Definition of Done

### Functional Completeness
- [ ] `ccexplorer get gcp` command returns accurate billing data
- [ ] All output formats (stdout, CSV, chart, Pinecone) work with GCP data
- [ ] Date filtering and granularity options function correctly
- [ ] GCP-specific filtering (project, service, region) works as expected
- [ ] Multi-project support enabled for enterprise users

### Code Quality
- [ ] All new code follows existing architectural patterns
- [ ] Unit test coverage ≥90% for new GCP functionality
- [ ] Integration tests pass with mock GCP environments
- [ ] Code passes all existing lint and security checks
- [ ] Generic utilities leveraged where appropriate

### Security & Compliance
- [ ] Authentication follows GCP security best practices
- [ ] No credentials stored in version control
- [ ] Proper error handling without information disclosure
- [ ] Security review completed and approved

### Documentation & Usability
- [ ] README updated with GCP setup instructions
- [ ] CLI help text provides clear usage guidance
- [ ] Error messages are user-friendly and actionable
- [ ] Code documentation follows existing standards

### Performance & Reliability
- [ ] GCP API calls implement proper rate limiting
- [ ] Error handling includes appropriate retry logic
- [ ] Performance benchmarks meet AWS equivalent metrics
- [ ] Memory usage patterns align with existing codebase

## Integration Points & Compatibility

### Backward Compatibility
- Maintain full backward compatibility with existing AWS functionality
- No changes to existing CLI command structure
- Preserve all current output formats and options

### Multi-Cloud Scenarios
- Consider future unified multi-cloud commands (`ccexplorer get combined`)
- Design data structures to support cross-cloud cost analysis
- Plan for unified reporting and visualization capabilities

### Extension Points
- Design GCP integration to support future cloud providers (Azure, etc.)
- Create reusable patterns for cloud provider authentication
- Build extensible filtering and aggregation frameworks

## Risk Assessment & Mitigation

### Technical Risks
- **GCP API Complexity**: Cloud Billing API has different data models than AWS Cost Explorer
  - *Mitigation*: Extensive prototyping and data mapping validation
- **Authentication Complexity**: GCP authentication differs from AWS credential chain
  - *Mitigation*: Leverage well-tested GCP SDK authentication patterns
- **Performance Impact**: Additional dependencies may affect build/runtime performance
  - *Mitigation*: Benchmark and optimize critical paths

### Business Risks  
- **Feature Scope Creep**: GCP has extensive billing features beyond core cost analysis
  - *Mitigation*: Focus on core cost analysis features first, plan additional features separately
- **Maintenance Overhead**: Supporting multiple cloud providers increases complexity
  - *Mitigation*: Strong architecture patterns and comprehensive testing strategy

## Success Metrics

- **Functionality**: Successfully retrieve and display GCP billing data equivalent to AWS capabilities
- **Performance**: GCP queries perform within 10% of AWS query times
- **Adoption**: Enable multi-cloud cost analysis for organizations using both AWS and GCP
- **Code Quality**: Maintain existing code quality standards while adding significant functionality

---

This issue represents a significant enhancement that positions ccExplorer as a comprehensive multi-cloud cost analysis tool. The implementation leverages existing architectural strengths while extending capabilities to support GCP's distinct billing ecosystem.

**Estimated Development Effort**: 6-8 weeks  
**Team Size**: 2-3 senior developers  
**Dependencies**: GCP billing API access, test environments