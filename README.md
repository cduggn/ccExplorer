
<h1 align="center"><code>ccExplorer</code></h1>

`ccExplorer` (Cloud cost explorer) is a simple command line tool to explore the 
cost of your cloud resources. It's not 
designed as a replacement for the official AWS CLI and does not offer the 
same exhaustive search option. He does however return results in a more
human-readable format, and orders them by cost in descending order.
It's primary use case is to surface costs based on pre-defined cost allocation tags. 



Quick Start
-----------

Once `ccExplorer` is installed you can run the help command to see the 
available commands.

```sh
ccexplorer --help
```
When you invoke a command, `ccExplorer` will use the AWS 
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
ccexplorer get aws -d SERVICE -t BucketName -u SERVICE="Amazon Simple 
Storage Service"
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

From source
    
```sh
git clone https://github.com/cduggn/ccExplorer.git
cd ccExplorer
go run .\cmd\ccexplorer\ccexplorer.go get aws -d SERVICE -d OPERATION -u SERVICE="Amazon DynamoDB"  -c
```


## Considerations when using Cost Explorer

- Cost Explorer accesses data for the last 12 months.
- Cost Explorer charges per paginated request.
- The AWS SDK uses the default credentials provider chain.
- Credits and refunds are automatically applied to Cost Explorer results. e`.
- Cost Explorer API calls can be tracked using CloudTrail. Requests are issued against us-east-1.
