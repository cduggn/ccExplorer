# cloudcost

A `Go` command line tool which supports exploring AWS costs using the AWS Cost Explorer SDK.

## Prerequisites

The SDK uses the AWS credential chain to find AWS credentials. The SDK looks for credentials in the following order: environment variables, shared credentials file, and EC2 instance profile. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Installation

See the releases page for the latest version.

## Commands

### `cost`
Specify granularity, and 'group by' parameters to get cost data. For example:

    $ cloudcost cost --granularity MONTHLY --group-by-dimension SERVICE --group-by-dimension AZ 
