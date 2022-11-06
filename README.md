# cloudcost

A `Go` command line tool which supports exploring AWS costs using the AWS Cost Explorer SDK. 

**Note**
The Cost Explorer API can access data for the last 13 months. This tool will only show data for the last 12 months.


## Prerequisites

The SDK uses the AWS credential chain to find AWS credentials. The SDK looks for credentials in the following order: environment variables, shared credentials file, and EC2 instance profile. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Installation

See the releases page for the latest version.

## Commands

**Note**
AWS billing information is updated up to three times daily. Querying the Cost Explorer is charged per paginated request. The default page size is 20. The default page size can be changed using the `--page-size` flag.

Cost Explorer supports the following commands:

### `cost`
Specify granularity, and 'group by' parameters to get cost data. For example:

    $ cloudcost cost --granularity MONTHLY --group-by-dimension SERVICE --group-by-dimension AZ 
