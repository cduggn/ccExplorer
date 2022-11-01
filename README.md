# cloudcost

A `Go` command line tool which supports exploring AWS costs using the AWS Cost Explorer SDK.

## Prerequisites

The SDK uses the AWS credential chain to find AWS credentials. The SDK looks for credentials in the following order: environment variables, shared credentials file, and EC2 instance profile. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Installation

Get the latest version of cloudcost library:
``` 
(In progress, not yet available)
go get github.com/cduggn/cloudcost
```

## Usage

The quickest way to get started is to run the following command:

```shell
$ cloudcost billing get "Amazon Route 53"
```