# cloudcost

`cloudcost` is a simple command line tool to track the cost of your cloud resources.
It is designed to be used with AWS, but could be extended to other cloud providers. It's primary 
use case is to surface costs based on pre-defined 
cost allocation tags. 
This approach simplifies the process of tracking costs across multiple projects and teams.   

**Note**
There are a number of considerations that 
The Cost Explorer API can access data for the last 12 months. This tool will only show data for the last 12 months.


## Considerations
 
There are a number of considerations that need to be taken into account before using this tool. 

- You will want to decide what cost allocation tags you will use to track costs. You will also need to ensure that the 
cost allocation tags are applied to all resources that you want to track. See [cost allocation tags.](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/cost-alloc-tags.html)
- The Cost Explorer API can access data for the last 12 months. This tool will only show data for the last 12 months.
- Cost Explorer charges per paginated request. Using cost allocation tags help reduce the number of requests that need to be made.
- The AWS SDK uses the default credentials provider chain. The SDK looks for credentials in the following order: environment variables, shared credentials file, and EC2 instance profile. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Installation

Precompiled binaries are available for Linux, Mac, and Windows on the releases [page](https://github.com/cduggn/cloudcost/releases).

## Commands

Cost Explorer supports the following commands:

### `get`

The `get` command returns the results of a Cost Explorer query. The `get` command supports the following sub flags:

`aws` - The AWS service to query. This is currently the only cloud service provider that can be queried.

The `aws` command provides the following sub flags:

Arguments:
   - --start-date string   The start date of the time period. The default is the current month.
   - --end-date string     The end date of the time period. The default is the current month.
   - --filter-by string    The filter to apply to the cost. The default is no filter. Used when the --group-by-tage flag is set.
   - --granularity string  The granularity of the cost. The default is DAILY.
   - --group-by-dimension string   The dimension to group the cost by. The default is [ SERVICE,USAGE_TYPE].
   - --group-by-tag string         The tag to group the cost by. The default is no grouping.

Basic usage:

The minimum required command is `get aws`. This will return the cost for the past 30 days. The default granularity is DAILY. The default group by dimension is [ SERVICE,USAGE_TYPE]. The default group by tag is no grouping.

    $ cloudcost get aws

Command returns the cost for the past 30 days grouped by the tag `ApplicationName` and the dimension `SERVICE`.
    
    $ cloudcost get aws --group-by-tag ApplicationName --group-by-dimension SERVICE

Command returns the cost for the past 30 days grouped by the tag `ApplicationName` and the dimension `SERVICE` and filter by the tag `ApplicationName` and the value `myapp`.
    
    $ cloudcost get aws --group-by-tag ApplicationName --filter-by myapp --group-by-dimension SERVICE

Command groups the cost by the dimension LINKED_ACCOUNT and filter by the tag `ApplicationName` and the value `myapp`.
    
    $ cloudcost get aws --group-by-dimension LINKED_ACCOUNT --group-by-tag ApplicationName--filter-by myapp

Command returns the cost for the past x days based on the provided start date. Refunds and credits are not filtered. UNBLENDED_COST cost is returned.

    $ cloudcost get aws  --group-by-tag ApplicationName --group-by-dimension SERVICE -r UNBLENDED_COST -g MONTHLY -s "2022-10-01"

Command returns the cost for the past x days based on the provided start date. Refunds and credits are filtered . UNBLENDED_COST costs are returned.

    $ cloudcost get aws  --group-by-tag ApplicationName --group-by-dimension SERVICE -r UNBLENDED_COST -g MONTHLY -s "2022-10-01" -c

Command returns the cost for the past x days based on the provided start date and groups by cost allocation tag filtered by specific value. Refunds and credits are filtered.

    $ cloudcost get aws  --group-by-tag ApplicationName --group-by-dimension SERVICE -r UNBLENDED_COST -g MONTHLY -s "2022-10-01" -c -f myapp

Dimension values include the following: AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, USAGE_TYPE_GROUP, RECORD_TYPE, and OPERATING_SYSTEM. For more information, see [Grouping and Filtering](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/billing-reports-costexplorer.html#ce-grouping-filtering).