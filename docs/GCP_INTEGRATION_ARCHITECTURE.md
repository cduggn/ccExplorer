# GCP Cloud Billing API Integration Architecture

## Overview

This document outlines the comprehensive architecture for integrating Google Cloud Platform (GCP) Cloud Billing API into ccExplorer, providing enterprise-grade multi-cloud cost analysis capabilities alongside the existing AWS Cost Explorer integration.

## Architecture Components

### 1. Service Layer Architecture (`internal/gcpservice/`)

#### Core Components

**`client.go`** - Main service implementation
- Implements `ports.GCPService` interface
- Enterprise-grade connection management with gRPC keepalives
- Built-in rate limiting (300 requests/minute for GCP Billing API)
- Exponential backoff retry logic with jitter
- Graceful error handling and resource cleanup

**`config.go`** - Authentication and configuration management
- Multiple authentication methods: Service Account, ADC
- Environment variable support (`GOOGLE_APPLICATION_CREDENTIALS`, `GCP_PROJECT_ID`)
- Multi-project enterprise support
- Organization and folder-level queries
- Secure credential handling

**`get_billing_data.go`** - Core billing data retrieval
- Service-level cost analysis with SKU breakdown
- Regional cost distribution
- Usage metrics extraction
- Credits and discounts processing
- Time-series data with configurable granularity

**`transformers.go`** - Data transformation pipeline
- Converts GCP billing data to internal types
- Enables reuse of existing output writers (stdout, CSV, charts, Pinecone)
- Aggregation and filtering capabilities
- Vector database preparation for AI/ML workflows

### 2. Interface Design (`internal/ports/services.go`)

The `GCPService` interface follows the same patterns as `AWSService` but adapted for GCP's distinct billing ecosystem:

```go
type GCPService interface {
    GetBillingData(ctx context.Context, req types.GCPBillingRequest) (*types.GCPBillingResponse, error)
    GetProjectBillingInfo(ctx context.Context, projectID string) (*billing.ProjectBillingInfo, error)
    GetBillingServices(ctx context.Context) (*billing.ListServicesResponse, error)
    GetBillingSkus(ctx context.Context, serviceID string) (*billing.ListSkusResponse, error)
    GetAssetInventory(ctx context.Context, req types.GCPAssetRequest) (*cloudasset.SearchAllResourcesResponse, error)
    GetCostForecast(ctx context.Context, req types.GCPForecastRequest) (*types.GCPForecastResponse, error)
}
```

### 3. Type System Extensions (`internal/types/command_types.go`)

Comprehensive GCP-specific types supporting:
- Multi-project queries (`GCPBillingRequest`)
- Enterprise hierarchy (Organizations, Folders)
- Detailed service breakdown (`GCPBillingItem`, `GCPSKU`)
- Regional cost analysis
- Credits and discounts tracking
- Time-series forecasting (`GCPForecastResponse`)

### 4. CLI Integration (`cmd/cli/get_command.go`)

**Command Structure:**
```bash
ccexplorer get gcp [flags]
ccexplorer get gcp forecast [flags]
```

**Key Features:**
- Project-level and organization-level queries
- Service and regional filtering
- Multiple output formats (stdout, CSV, charts, Pinecone)
- Cost threshold filtering
- Credits/discounts inclusion control

## Security Architecture

### 1. Authentication Security

**Service Account Security:**
- Secure key file handling with proper file permissions
- No credentials stored in version control
- Support for key rotation through standard GCP mechanisms
- Environment-based configuration (`GOOGLE_APPLICATION_CREDENTIALS`)

**IAM Requirements:**
- Minimum required permissions: `billing.viewer`, `cloudasset.viewer`
- Optional: `resourcemanager.viewer` for organization queries
- Principle of least privilege enforcement

**Application Default Credentials (ADC):**
- Fallback authentication method
- Supports workload identity for GKE deployments
- Local development support through `gcloud auth application-default login`

### 2. API Security

**Rate Limiting:**
- Token bucket implementation (5 requests/second, 10 burst capacity)
- Respects GCP Billing API limits (300 requests/minute)
- Configurable limits for different API endpoints

**Request Validation:**
- Input sanitization for all user-provided data
- Project ID format validation
- Date range validation
- Parameter bounds checking

**Error Handling:**
- Sanitized error messages (no credential exposure)
- Structured error logging for security monitoring
- Proper HTTP status code handling

### 3. Data Protection

**Data Encryption:**
- TLS 1.3 for all API communications
- Certificate validation enforced
- No plaintext credential storage

**PII Handling:**
- Resource names and labels may contain PII
- Configurable data masking for sensitive outputs
- Audit logging for data access

**Access Control:**
- Verify user has appropriate GCP billing permissions
- Project-level access control validation
- Organization-level permission checks

## Performance Architecture

### 1. Rate Limiting Strategy

**Implementation:**
- Token bucket algorithm with configurable parameters
- Per-service rate limiting (billing API vs asset API)
- Backoff strategy for quota exceeded errors

**Configuration:**
```go
rateLimiter := &RateLimiter{
    tokensPerSecond: 5.0,    // 300 requests / 60 seconds
    maxTokens:       10,     // Burst capacity
    tokens:          10,
    lastRefill:      time.Now(),
}
```

### 2. Connection Management

**gRPC Optimization:**
- Connection pooling with keepalive parameters
- Maximum message size configuration (100MB)
- Connection reuse across requests
- Graceful shutdown handling

**Retry Logic:**
- Exponential backoff with jitter
- Maximum 3 retries with configurable delays
- Retryable error detection (5xx, timeout, unavailable)

### 3. Data Processing Optimization

**Batch Processing:**
- Multi-project queries in parallel (configurable concurrency)
- SKU data caching for repeated queries
- Pagination handling for large result sets

**Memory Management:**
- Streaming response processing
- Configurable result set limits
- Efficient data transformation pipelines

## Testing Strategy

### 1. Unit Testing

**Service Layer Tests:**
```go
// internal/gcpservice/client_test.go
func TestGCPService_Authentication(t *testing.T)
func TestGCPService_RateLimiting(t *testing.T)
func TestGCPService_RetryLogic(t *testing.T)
func TestGCPService_ErrorHandling(t *testing.T)
```

**Transformation Tests:**
```go
// internal/gcpservice/transformers_test.go
func TestTransformToServices(t *testing.T)
func TestAggregateData(t *testing.T)
func TestApplyFiltering(t *testing.T)
```

**Type Validation Tests:**
```go
// internal/types/command_types_test.go
func TestGCPBillingRequest_Validate(t *testing.T)
func TestGCPBillingRequest_Equals(t *testing.T)
```

### 2. Integration Testing

**API Integration:**
- Mock GCP API responses using `httptest`
- Test authentication flow with dummy credentials
- Validate API request formatting and response parsing

**CLI Testing:**
```go
// cmd/cli/gcp_commands_test.go
func TestGCPCommand_FlagParsing(t *testing.T)
func TestGCPCommand_InputValidation(t *testing.T)
func TestGCPCommand_OutputFormatting(t *testing.T)
```

### 3. End-to-End Testing

**Workflow Tests:**
- Complete command execution with mock APIs
- Output writer integration tests
- Multi-format output validation

**Error Scenario Tests:**
- Invalid credentials handling
- Network failure recovery
- API quota exceeded scenarios

### 4. Performance Testing

**Load Testing:**
- Concurrent request handling
- Memory usage under high load
- Rate limiting effectiveness

**Benchmarking:**
```go
func BenchmarkGCPService_GetBillingData(b *testing.B)
func BenchmarkDataTransformation(b *testing.B)
```

## Data Flow Architecture

### 1. Request Flow

```
CLI Input → GCPCommandType → GCPBillingRequest → GCPService → GCP APIs
```

1. **CLI Processing:** Flag parsing and validation
2. **Request Synthesis:** Convert CLI input to internal request types
3. **Service Layer:** Rate limiting, authentication, API calls
4. **Response Processing:** Error handling, data transformation

### 2. Response Flow

```
GCP APIs → GCPBillingResponse → DataTransformer → Writer → Output
```

1. **API Response:** Raw GCP billing data
2. **Transformation:** Convert to internal types
3. **Writer Selection:** Based on output format preference
4. **Output Generation:** stdout, CSV, charts, or Pinecone

### 3. Error Flow

```
Error → Classification → Retry Logic → User-Friendly Message
```

1. **Error Detection:** API errors, validation errors, network errors
2. **Classification:** Retryable vs non-retryable
3. **Retry Processing:** Exponential backoff for transient errors
4. **User Communication:** Clear, actionable error messages

## Implementation Phases

### Phase 1: Foundation (Completed)
- ✅ Service interface design
- ✅ Basic client implementation
- ✅ Authentication configuration
- ✅ Type system extensions

### Phase 2: Core Integration (Completed)
- ✅ Billing API integration
- ✅ CLI command structure
- ✅ Data transformation pipeline
- ✅ Error handling framework

### Phase 3: Advanced Features (Ready for Implementation)
- 🔄 Asset inventory integration
- 🔄 Cost forecasting algorithms
- 🔄 Advanced filtering and aggregation
- 🔄 Performance optimizations

### Phase 4: Testing & Documentation (In Progress)
- 🔄 Comprehensive unit tests
- 🔄 Integration test suite
- ✅ Security review
- ✅ Performance benchmarking

## Deployment Considerations

### 1. Dependencies

**Required GCP APIs:**
- Cloud Billing API
- Cloud Asset Inventory API (optional)
- Cloud Resource Manager API (for organization queries)

**Go Dependencies:**
```go
cloud.google.com/go/billing v1.18.0
cloud.google.com/go/asset v1.19.0
google.golang.org/api v0.155.0
golang.org/x/oauth2 v0.15.0
```

### 2. Configuration

**Environment Variables:**
```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
export GCP_PROJECT_ID="my-project-id"
export GCP_BILLING_ACCOUNT_ID="012345-678901-234567"
export GCP_ORGANIZATION_ID="123456789"
```

**IAM Setup:**
```bash
# Create service account
gcloud iam service-accounts create ccexplorer-billing

# Grant required permissions
gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:ccexplorer-billing@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/billing.viewer"

# Generate and download key
gcloud iam service-accounts keys create ccexplorer-key.json \
    --iam-account=ccexplorer-billing@PROJECT_ID.iam.gserviceaccount.com
```

### 3. Monitoring and Observability

**Metrics:**
- API request latency and success rates
- Rate limiting statistics
- Cost query volumes and patterns
- Error rates by type

**Logging:**
- Structured logging with correlation IDs
- Security event logging
- Performance metrics
- User activity tracking

## Future Enhancements

### 1. Advanced Analytics
- Cost anomaly detection using historical patterns
- Budget alerting integration
- Cost optimization recommendations
- Multi-cloud cost comparison dashboards

### 2. Enterprise Features
- RBAC integration with Google Cloud Identity
- Cost allocation by teams/departments
- Automated reporting and scheduling
- Integration with FinOps workflows

### 3. Performance Improvements
- Response caching for frequently accessed data
- Incremental data updates
- Parallel processing for large organizations
- Real-time streaming for cost alerts

## Conclusion

This GCP integration architecture provides enterprise-grade multi-cloud cost analysis capabilities while maintaining the clean architecture principles of ccExplorer. The design ensures scalability, security, and maintainability while enabling seamless integration with existing AWS functionality.

The implementation follows industry best practices for cloud API integration, provides comprehensive error handling and security features, and maintains compatibility with all existing output formats and workflows.