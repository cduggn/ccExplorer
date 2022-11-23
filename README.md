# cloudcost

`cloudcost` is a simple command line tool to track the cost of your cloud resources.
It is designed to be used with AWS, but could be extended to other cloud providers. It's primary 
use case is to surface costs based on pre-defined 
cost allocation tags. 
This approach simplifies the process of tracking costs across multiple projects and teams.   


## Considerations
 
There are a number of considerations that need to be taken into account before using this tool. 

- You will want to decide what cost allocation tags you will use to track costs. You will also need to ensure that the 
cost allocation tags are applied to all resources that you want to track. See [cost allocation tags.](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/cost-alloc-tags.html)
- The Cost Explorer API can access data for the last 12 months. This tool will only show data for the last 12 months.
- Cost Explorer charges per paginated request. Using cost allocation tags help reduce the number of requests that need to be made.
- The AWS SDK uses the default credentials provider chain. The SDK looks for credentials in the following order: environment variables, 
shared credentials file, and EC2 instance profile or ECS task definition if running on either platform. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).
- Unblended costs are used to calculate the total cost of a resource. Unblended costs are the sum of the costs of all usage of a resource. This is the default cost tyoe returned by the tool. 
- Credits and refunds are automatically applied to your account. Both can be excluded from the cost data by setting the `exclude_credit` flag to `true`.

## Installation

Precompiled binaries are available for Linux, Mac, and Windows on the releases [page](https://github.com/cduggn/cloudcost/releases).

## Usage

Cost Explorer supports the `get aws` command and subcommand with the following options:

```bash
GetBill = DESCRIPTION
Fetches billing information for the time interval provided using the AWS Cost Explorer API

Prerequisites:
- AWS credentials must be configured in ~/.aws/credentials
- AWS region must be configured in ~/.aws/config
- Cost Allocation Tags if you want to filter by tag ( Note cost allocation tags can take up to 24 hours to be applied )

Usage:
  cloudcost get aws [flags]

Flags:
  -e, --end-date string              End date for billing information. Default is todays date. (default "2022-11-23")
  -c, --exclude-credit               Exclude credit and refund information in the report. This is enabled by default
  -f, --filter-by string             When grouping by tag, filter by tag value
  -g, --granularity string           Granularity of billing information to fetch (default "DAILY")
  -d, --group-by-dimension strings   Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, USAGE_TYPE ] (default [SERVICE,USAGE_TYPE])
  -t, --group-by-tag string          Group by cost allocation tag
  -h, --help                         help for aws
  -r, --rates strings                Cost and Usage rates to fetch [ Rates: BLENDED_COST, UNBLENDED_COST, AMORTIZED_COST, NET_AMORTIZED_COST, NET_UNBLENDED_COST, USAGE_QUANTITY ]. Defaults to UNBLENDED_COST (default [UNBLENDED_COST])
  -s, --start-date string            Start date for billing information. Defaults to the past 7 days (default "2022-10-24")
```

## Examples
A sample command to fetch billing information for the past 7 days grouped by service and usage type:

```bash
cloudcost get aws -s 2021-10-24 -e 2021-10-31 -g DAILY -r UNBLENDED_COST -g SERVICE -g USAGE_TYPE
```

Using cost allocation tags to filter by project:

```bash
cloudcost get aws -t Project -f my-project -d SERVICE -g MONTHLY -s 2022-10-01 -e 2022-11-21 -r UNBLENDED_COST
```

