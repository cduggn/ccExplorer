
<h1 align="center"><code>ccExplorer</code></h1>

`ccExplorer` (Cloud cost explorer) is a simple command line tool to explore the 
cost of your cloud resources. It's not 
designed as a replacement for the official AWS CLI. Instead of returning 
results as JSON, it returns a human readable table with costs displayed in 
descending order by cost.
It is designed to be used with AWS, but could be extended to other cloud providers. It's primary 
use case is to surface costs based on pre-defined 
cost allocation tags. 
This approach simplifies the process of tracking costs across multiple projects and teams.   



Quick Start
-----------

Once `ccExplorer` is installed you can run the help command to see the 
available commands.

```sh
ccexplorer --help
```
When you invoke a command, `ccExplorer` will follow use the AWS 
credential chain to authenticate with AWS.

If no cost allocation tags have been defined, the  `ccExplorer` can still be 
used to 
filter and group resources based on their 
AWS resource types. This can be achieved by using the group by and filter 
flags 

```sh
ccexplorer get aws -d SERVICE -d OPERATION c
```

This will return a list of all AWS services and operations that have been
used in the specified billing period. If no billing period is specified, the
current calendar month will be used. Results are sorted by cost in 
descending order and refunds, discounts and credits are excluded.

```sh
ccexplorer get aws -t ApplicationName -d OPERATION -s 2022-12-10 -f "my-project"
```
This will return a list of all AWS operations that have been used in the 
specified billing period for the specified project. The `-f` flag can be
used to filter results based on the value of a cost allocation tag. If no 
filter value is specified, all resources will be returned. Results are  
sorted by cost in descending order.

```sh
ccexplorer get aws -d SERVICE -t ApplicationName -u SERVICE="Amazon Simple Storage Service"  -c
```

This will return a list of costs for S3 buckets that have been used in the
specified billing period and that have been tagged with the ApplicationName
tag. Results are sorted by cost in descending order.

```sh
ccexplorer get aws -d SERVICE -t ApplicationName -u SERVICE="Amazon Simple Storage Service"  -c -f "my-application"
```

This will return a list of costs for the specified application that have
been used in the specified billing period. Results are sorted by cost in
descending order.

```sh
ccexplorer get aws -d SERVICE -t BucketName -u SERVICE="Amazon Simple Storage"
```

This will return a list of costs for S3 buckets filtered by the bucket name
tag. Results are sorted by cost in descending order.


Cost Forecast
-------------

The cost forecast command supports both wide ranging and granular forecasts.

```sh 
ccexplorer get aws forecast -e 2023-01-21 -d SERVICE="AWS Lambda"
```

This will return a forecast for the cost of AWS Lambda for the current 
billing period and the next 12 months. The forecast is based on the current
usage of AWS Lambda and the average cost of AWS Lambda over the last 12.


```sh 
ccexplorer get aws -d OPERATION -t ApplicationName -u OPERATION="PutObject"  -c
```

This will return a list of costs grouped by application name for the
PutObject operation. Results are sorted by cost in descending order.


Installation
------------
Precompiled binaries are available for Linux, Mac, and Windows on the releases [page](https://github.com/cduggn/cloudcost/releases).


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
