# GCP Cloud Billing API Integration - Implementation Guide

## Quick Start

### 1. Add Dependencies

Update `go.mod` with GCP SDK dependencies:

```bash
go get cloud.google.com/go/billing@v1.18.0
go get cloud.google.com/go/asset@v1.19.0  
go get google.golang.org/api@v0.155.0
go mod tidy
```

### 2. Set up GCP Authentication

**Option A: Service Account (Recommended for Production)**
```bash
# Create service account
gcloud iam service-accounts create ccexplorer-billing \
    --description="ccExplorer GCP Billing Access" \
    --display-name="ccExplorer Billing"

# Add billing viewer role
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:ccexplorer-billing@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/billing.viewer"

# Generate key file
gcloud iam service-accounts keys create ~/.config/gcp/ccexplorer-key.json \
    --iam-account=ccexplorer-billing@YOUR_PROJECT_ID.iam.gserviceaccount.com

# Set environment variable
export GOOGLE_APPLICATION_CREDENTIALS="$HOME/.config/gcp/ccexplorer-key.json"
export GCP_PROJECT_ID="YOUR_PROJECT_ID"
```

**Option B: Application Default Credentials (Development)**
```bash
gcloud auth application-default login
export GCP_PROJECT_ID="YOUR_PROJECT_ID"
```

### 3. Test GCP Integration

```bash
# Build with GCP support
make build

# Test basic GCP cost query
./bin/ccexplorer get gcp --project YOUR_PROJECT_ID

# Test with filtering
./bin/ccexplorer get gcp --project YOUR_PROJECT_ID --services "Compute Engine" --granularity DAILY

# Test CSV export
./bin/ccexplorer get gcp --project YOUR_PROJECT_ID --printFormat csv
```

## Complete Implementation Details

### File Structure Overview

```
ccExplorer/
├── internal/
│   ├── ports/
│   │   └── services.go                 # Extended with GCPService interface
│   ├── gcpservice/                     # New GCP service package
│   │   ├── client.go                   # Main GCP service implementation
│   │   ├── config.go                   # Authentication & configuration
│   │   ├── get_billing_data.go         # Billing data retrieval
│   │   └── transformers.go             # Data transformation pipeline
│   └── types/
│       └── command_types.go            # Extended with GCP types
├── cmd/cli/
│   └── get_command.go                  # Extended with GCP commands
└── docs/
    ├── GCP_INTEGRATION_ARCHITECTURE.md # Architecture documentation
    └── GCP_IMPLEMENTATION_GUIDE.md     # This file
```

### Core Implementation Components

#### 1. Service Interface (`internal/ports/services.go`)

The `GCPService` interface provides comprehensive billing capabilities:

```go
type GCPService interface {
    // Core billing data retrieval
    GetBillingData(ctx context.Context, req types.GCPBillingRequest) (*types.GCPBillingResponse, error)
    
    // Project and organizational billing info
    GetProjectBillingInfo(ctx context.Context, projectID string) (*billing.ProjectBillingInfo, error)
    
    // Service discovery and SKU analysis
    GetBillingServices(ctx context.Context) (*billing.ListServicesResponse, error)
    GetBillingSkus(ctx context.Context, serviceID string) (*billing.ListSkusResponse, error)
    
    // Asset inventory integration
    GetAssetInventory(ctx context.Context, req types.GCPAssetRequest) (*cloudasset.SearchAllResourcesResponse, error)
    
    // Cost forecasting
    GetCostForecast(ctx context.Context, req types.GCPForecastRequest) (*types.GCPForecastResponse, error)
}
```

#### 2. Client Implementation (`internal/gcpservice/client.go`)

Enterprise-grade features:
- **Rate Limiting**: Token bucket with 5 req/sec, 10 burst capacity
- **Retry Logic**: Exponential backoff with jitter for transient failures
- **Connection Management**: gRPC keepalives and connection pooling
- **Error Handling**: Comprehensive error classification and recovery

Key Methods:
```go
func New(opts ...option.ClientOption) (*Service, error)
func (s *Service) GetBillingData(ctx context.Context, req types.GCPBillingRequest) (*types.GCPBillingResponse, error)
func (s *Service) Close() error
```

#### 3. Configuration Management (`internal/gcpservice/config.go`)

Multi-environment configuration support:
```go
type Config struct {
    // Authentication
    ServiceAccountPath string
    ProjectID          string
    BillingAccountID   string
    
    // Enterprise features
    OrganizationID     string
    FolderIDs          []string
    ProjectIDs         []string
    
    // Performance settings
    MaxConcurrentCalls int
    RequestTimeout     int
    DefaultCurrency    string
}
```

Environment variable support:
- `GOOGLE_APPLICATION_CREDENTIALS`
- `GCP_PROJECT_ID`
- `GCP_BILLING_ACCOUNT_ID`
- `GCP_ORGANIZATION_ID`

#### 4. Data Transformation (`internal/gcpservice/transformers.go`)

Transforms GCP data to work with existing ccExplorer infrastructure:

```go
type GCPDataTransformer struct {
    config *Config
}

// Key transformation methods
func (t *GCPDataTransformer) TransformToServices(response *types.GCPBillingResponse) (map[int]types.Service, error)
func (t *GCPDataTransformer) TransformToCSVData(response *types.GCPBillingResponse) ([][]string, error)
func (t *GCPDataTransformer) TransformToVectorData(response *types.GCPBillingResponse) ([]types.VectorStoreItem, error)
```

#### 5. Type System (`internal/types/command_types.go`)

Comprehensive GCP-specific types:

**Request Types:**
```go
type GCPBillingRequest struct {
    ProjectID          string
    ProjectIDs         []string
    BillingAccount     string
    OrganizationID     string
    Services           []string
    Regions            []string
    Time               GCPTimeRange
    Granularity        string
    Currency           string
    // ... additional fields
}
```

**Response Types:**
```go
type GCPBillingResponse struct {
    Services        map[string]GCPBillingItem
    ProjectInfo     string
    BillingAccount  string
    TotalCost       float64
    TotalCredits    float64
    Metadata        GCPResponseMetadata
}
```

#### 6. CLI Integration (`cmd/cli/get_command.go`)

Extended command structure:
```bash
ccexplorer get gcp [flags]                    # Main GCP cost analysis
ccexplorer get gcp forecast [flags]           # GCP cost forecasting
```

**Key Flags:**
- `--project` (required): GCP Project ID
- `--projects`: Multiple project IDs
- `--services`: Filter by service names
- `--regions`: Filter by regions
- `--granularity`: DAILY, MONTHLY, HOURLY
- `--printFormat`: stdout, csv, chart, pinecone

## Advanced Usage Examples

### Multi-Project Analysis
```bash
# Analyze multiple projects together
ccexplorer get gcp --project main-project \
  --projects "dev-project,staging-project,prod-project" \
  --granularity MONTHLY

# Organization-wide analysis
ccexplorer get gcp --organization 123456789 --granularity MONTHLY
```

### Service-Specific Analysis
```bash
# Compute Engine costs only
ccexplorer get gcp --project my-project \
  --services "Compute Engine" \
  --regions "us-central1,us-east1"

# Multiple services with regional breakdown
ccexplorer get gcp --project my-project \
  --services "Compute Engine,Cloud Storage,BigQuery" \
  --groupBy service,region
```

### Cost Optimization Queries
```bash
# Find services costing more than $100
ccexplorer get gcp --project my-project \
  --costThreshold 100.00 \
  --sortByDate=false

# Exclude credits to see true usage costs
ccexplorer get gcp --project my-project \
  --includeCredits=false \
  --includeDiscounts=false
```

### Export and Analysis
```bash
# Export to CSV for spreadsheet analysis
ccexplorer get gcp --project my-project \
  --printFormat csv \
  --startDate 2024-01-01 \
  --endDate 2024-12-31

# Generate charts for visualization
ccexplorer get gcp --project my-project \
  --printFormat chart \
  --granularity DAILY

# Store in vector database for AI analysis
ccexplorer get gcp --project my-project \
  --printFormat pinecone
```

### Forecasting
```bash
# Basic cost forecast
ccexplorer get gcp forecast --project my-project

# Service-specific forecast
ccexplorer get gcp forecast --project my-project \
  --services "Compute Engine,BigQuery" \
  --granularity DAILY
```

## Security Best Practices

### 1. Credential Management
```bash
# Secure service account key storage
sudo mkdir -p /etc/ccexplorer/credentials
sudo chown root:ccexplorer-group /etc/ccexplorer/credentials
sudo chmod 750 /etc/ccexplorer/credentials

# Set restrictive permissions on key files
chmod 600 /etc/ccexplorer/credentials/service-account.json
```

### 2. IAM Configuration
```bash
# Minimal required permissions
gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:ccexplorer@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/billing.viewer"

# For organization-level access
gcloud organizations add-iam-policy-binding ORG_ID \
    --member="serviceAccount:ccexplorer@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/billing.viewer"
```

### 3. Network Security
```bash
# Configure firewall rules for API access
gcloud compute firewall-rules create allow-gcp-apis \
    --allow tcp:443 \
    --source-ranges "10.0.0.0/8" \
    --description "Allow GCP API access"
```

## Troubleshooting

### Common Issues

**1. Authentication Errors**
```bash
# Verify credentials
gcloud auth application-default print-access-token

# Check service account permissions
gcloud projects get-iam-policy PROJECT_ID \
    --flatten="bindings[].members" \
    --format="table(bindings.role)" \
    --filter="bindings.members:*ccexplorer*"
```

**2. API Access Issues**
```bash
# Enable required APIs
gcloud services enable cloudbilling.googleapis.com
gcloud services enable cloudasset.googleapis.com
gcloud services enable cloudresourcemanager.googleapis.com
```

**3. Rate Limiting**
```bash
# Check quota usage
gcloud logging read "resource.type=gce_instance AND \
  protoPayload.methodName=cloudbilling.googleapis.com" \
  --limit=50 --format="table(timestamp,protoPayload.status.code)"
```

### Debug Mode
```bash
# Enable verbose logging
export CCEXPLORER_DEBUG=true
export GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE=true

# Run with debug output
./bin/ccexplorer get gcp --project my-project 2>&1 | tee debug.log
```

## Performance Tuning

### 1. Rate Limiting Configuration
```go
// Adjust rate limits in config
rateLimiter := &RateLimiter{
    tokensPerSecond: 10.0,  // Increase for higher quotas
    maxTokens:       20,    // Increase burst capacity
}
```

### 2. Concurrency Settings
```go
config := &Config{
    MaxConcurrentCalls: 5,   // Increase for better throughput
    RequestTimeout:     60,  // Increase for large queries
}
```

### 3. Caching Strategy
```bash
# Enable response caching (when implemented)
export GCP_CACHE_TTL=300  # 5 minutes
export GCP_CACHE_SIZE=100 # 100 MB
```

## Testing

### Unit Tests
```bash
# Run GCP-specific tests
go test ./internal/gcpservice/... -v

# Run with coverage
go test ./internal/gcpservice/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests
```bash
# Test with mock GCP APIs
go test ./internal/gcpservice/... -tags=integration

# Test CLI commands
go test ./cmd/cli/... -run=TestGCP
```

### Performance Benchmarks
```bash
# Run benchmarks
go test ./internal/gcpservice/... -bench=. -benchmem

# Profile memory usage
go test ./internal/gcpservice/... -bench=BenchmarkGetBillingData -memprofile=mem.prof
go tool pprof mem.prof
```

## Monitoring and Observability

### Metrics Collection
```bash
# Monitor API usage
curl -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  "https://monitoring.googleapis.com/v3/projects/PROJECT_ID/metricDescriptors"
```

### Logging Configuration
```bash
# Enable Cloud Logging
export GOOGLE_CLOUD_LOGGING=true

# Set log level
export CCEXPLORER_LOG_LEVEL=info
```

## Next Steps

1. **Enable APIs**: Ensure Cloud Billing API is enabled in your GCP project
2. **Set up Authentication**: Configure service account or ADC
3. **Install Dependencies**: Add GCP SDK dependencies to your project
4. **Test Integration**: Run basic queries to verify functionality
5. **Configure Monitoring**: Set up logging and metrics collection
6. **Scale Usage**: Gradually increase query volume and complexity

For additional support and advanced configuration options, refer to the comprehensive architecture documentation in `GCP_INTEGRATION_ARCHITECTURE.md`.