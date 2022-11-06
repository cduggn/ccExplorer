# cloudcost

A `Go` command line tool which supports exploring AWS costs using the AWS Cost Explorer SDK. 

**Note**
The Cost Explorer API can access data for the last 12 months. This tool will only show data for the last 12 months.


## Prerequisites

The SDK uses the AWS credential chain to find AWS credentials. The SDK looks for credentials in the following order: environment variables, shared credentials file, and EC2 instance profile. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Installation

See the releases page for the latest version.

## Commands

**Note**
AWS billing information is updated up to three times daily. Querying the Cost Explorer is charged per paginated request. The default page size is 20. The default page size can be changed using the `--page-size` flag.

Cost Explorer supports the following commands:

### `cost`

The `cost` command returns the total cost for the given time period. The time period can be specified using the `--start-date` and `--end-date` flags. The default time period is the current month.

Arguments:
   - --start-date string   The start date of the time period. The default is the current month.
   - --end-date string     The end date of the time period. The default is the current month.
   - --filter-by string    The filter to apply to the cost. The default is no filter. Used when the --group-by-tage flag is set.
   - --granularity string  The granularity of the cost. The default is DAILY.
   - --group-by-dimension string   The dimension to group the cost by. The default is [ SERVICE,USAGE_TYPE].
   - --group-by-tag string         The tag to group the cost by. The default is no grouping.


    $ cloudcost aws get-cost-and-usage -g DAILY -t <Cost-Allocation-Tag-Name> -d SERVICE -s 2022-10-01 -f "Some value to filter on"

Sample output when filtering on tag named ApplicationName

```bash
+-------------------------------------+-------------------------+---------------+-------------------+------+
| DIMENSION/TAG                       | DIMENSION/TAG           | METRIC NAME   | AMOUNT            | UNIT |
+-------------------------------------+-------------------------+---------------+-------------------+------+
| AWS Data Transfer                   | ApplicationName$        | BlendedCost   | 0.000000000       | USD  |
| AWS Data Transfer                   | ApplicationName$        | UnblendedCost | 0.000000000       | USD  |
| AWS Data Transfer                   | ApplicationName$        | UsageQuantity | 0                 | N/A  |
| AWS Step Functions                  | ApplicationName$        | UnblendedCost | 0                 | USD  |
| AWS Step Functions                  | ApplicationName$        | UsageQuantity | 0                 | N/A  |
| AWS Step Functions                  | ApplicationName$        | BlendedCost   | 0                 | USD  |
| Amazon API Gateway                  | ApplicationName$        | BlendedCost   | 0                 | USD  |
| Amazon API Gateway                  | ApplicationName$        | UnblendedCost | 0                 | USD  |
| Amazon API Gateway                  | ApplicationName$        | UsageQuantity | 0.0000000000      | N/A  |
| Amazon CloudFront                   | ApplicationName$        | BlendedCost   | 0.0000000000      | USD  |
| Amazon CloudFront                   | ApplicationName$        | UnblendedCost | 0.0000000000      | USD  |
| Amazon CloudFront                   | ApplicationName$        | UsageQuantity | 0.0000000000      | N/A  |
| Refund                              | ApplicationName$        | BlendedCost   | 0                 | USD  |
| Refund                              | ApplicationName$        | UnblendedCost | 0                 | USD  |

```