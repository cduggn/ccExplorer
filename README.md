# ccExplorer

`ccExplorer` (Cloud cost explorer) is a simple command line tool to explore the 
cost of your cloud resources and forecasts for future costs.
It is designed to be used with AWS, but could be extended to other cloud providers. It's primary 
use case is to surface costs based on pre-defined 
cost allocation tags. 
This approach simplifies the process of tracking costs across multiple projects and teams.   

## Installation

Precompiled binaries are available for Linux, Mac, and Windows on the releases [page](https://github.com/cduggn/cloudcost/releases).

## Usage

Cost Explorer supports the `get aws` command and subcommand with the following 
options:

```bash
Command: aws
Description: Returns cost and usage summary for the specified time period.

Prerequisites:
- AWS credentials must be configured in ~/.aws/credentials
- AWS region must be configured in ~/.aws/config
- Cost Allocation Tags must exist in AWS console if you want to filter by tag (
Note cost allocation tags can take up to 24 hours to be applied )

Usage:
  ccxplorer get aws [flags]
  ccxplorer get aws [command]

Available Commands:
  forecast    Return cost, usage, and resoucrce information including ARN

Flags:
  -d, --dimensions strings   Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, USAGE_TYPE ]
  -e, --end string           Defaults to the present day (default "2022-12-29")
  -c, --exclude-discounts    Exclude credit, refund, and discount information in the report summary. Disabled by default.
  -f, --filter-by string     When grouping by tag, filter by tag value
  -g, --granularity string   Granularity of billing information to fetch. Monthly, Daily or Hourly (default "MONTHLY")
  -h, --help                 help for aws
  -s, --start string         Defaults to the start of the current month (default "2022-12-01")
  -t, --tags string          Group by cost allocation tag
```

## Cost And Usage Report 
This command fetches the cost and usage information for the 
authenticated 
account. Results are grouped by SERVICE name and OPERATION type and excludes 
and credits or discounts. 

```bash
Note: The dimension types impact the volume of results returned. )

```bash) 
service 
and operation type:

```bash
cloudcost get aws -d SERVICE -d OPERATION c

```

Results are grouped by SERVICE name and OPERATION type in descending order 
by cost.

<style>
table {
  font-size: 8pt;
}
</style>

| RANK | DIMENSION/TAG   | DIMENSION/TAG   | METRIC NAME | NUMERIC AMOUNT | STRING AMOUNT | UNIT | GRANULARITY | START | END  |
|---------|-----------|--------|------|------| ------|------|------|------|------|
| 1 | Amazon Route 53   | HostedZone | UnblendedCost |  1.50000010 | 1.5  | USD | MONTHLY | 2021-12-01 | 2021-12-31 |
| 2 | AWS Cost Explorer  | GetCostAndUsage | UnblendedCost | 0.46000010  | 0.46 | USD | MONTHLY | 2021-12-01 | 2021-12-31 |
| 3 | Amazon Route 53  | Health-Check-HTTPS | UnblendedCost | 0.22580610|   0.23 | USD | MONTHLY | 2021-12-01 | 2021-12-31 |
| 4 | AWS Config   | None | UnblendedCost | 0.18900010 | 0.19 | USD | MONTHLY | 2021-12-01 | 2021-12-31 |
| 5 | Amazon Route 53   | Health-Check-HTTPS | UnblendedCost | 0.18900010 | 0.19 | USD | MONTHLY | 2021-12-01 | 2021-12-31 |

Using cost allocation tags to filter by project:

```bash
cloudcost get aws -t ApplicationName -d OPERATION -s 2022-12-10 -f "my-project"
```

## Cost Forecast
The cost forecast command supports both wide ranging and granular forecasts.

Example 1: Cost forecast for AWS Lambda given a specific end date:

```bash 
cloudcost get aws forecast -e 2023-01-21 -d SERVICE="AWS Lambda"
```

Example 2: Cost forecast for S3's PutObject API. 


```bash
cloudcost get aws forecast -e 2023-01-21 -d OPERATION="PutObject"
```


## Considerations when using Cost Explorer

There are a number of considerations that need to be taken into account before using this tool.

- You will want to decide what cost allocation tags you will use to track costs. You will also need to ensure that the
  cost allocation tags are applied to all resources that you want to track. See [cost allocation tags.](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/cost-alloc-tags.html)
- The Cost Explorer API can access data for the last 12 months. This tool will only show data for the last 12 months.
- Cost Explorer charges per paginated request. Using cost allocation tags help reduce the number of requests that need to be made.
- The AWS SDK uses the default credentials provider chain. The SDK looks for credentials in the following order: environment variables,
  shared credentials file, and EC2 instance profile or ECS task definition if running on either platform. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).
- Unblended costs are used to calculate the total cost of a resource. Unblended costs are the sum of the costs of all usage of a resource. This is the default cost tyoe returned by the tool.
- Credits and refunds are automatically applied to your account. Both can be excluded from the cost data by setting the `exclude_credit` flag to `true`.
- Cost Explorer API calls can be expensive. The tool will cache the results of the API calls to reduce the number of calls that need to be made. The cache is stored in the `~/.cloudcost` directory. [in-progress]
- Cost Explorer API calls can be tracked using CloudTrail. Requests are issued against us-east-1.
- By default CLI shows data from the beginning of the previous month
