###### ccExplorer is in Alpha

<h1 align="center"><code>ccExplorer</code></h1>

<hr>
<div align="center">
<a href="https://github.com/cduggn/ccExplorer/actions" 
alt="goreleaser status">
<img src="https://github.com/cduggn/ccExplorer/actions/workflows/release.yml/badge.svg">
</a>
<a href="https://goreportcard.com/report/github.com/cduggn/ccexplorer">
    <img src="https://goreportcard.com/badge/github.com/cduggn/ccexplorer" alt="Go Report Card">
</a>
<a href="https://github.com/cduggn/ccExplorer/actions" 
alt="CodeQL status">
<img src="https://github.com/cduggn/ccExplorer/actions/workflows/codeql.yml/badge.svg">
</a>
<a href="https://github.com/cduggn/ccExplorer/releases" 
alt="release status">
<img src="https://img.shields.io/github/v/release/cduggn/ccExplorer">
</a>


</div>

`ccExplorer` (Cloud cost explorer) is a simple command line tool to explore the 
cost of your cloud resources. It's built on opensource tools like [cobra](https://github.com/spf13/cobra),
[go-echarts](https://github.com/go-echarts/go-echarts), and [go-pretty](https://github.com/jedib0t/go-pretty).
It lets you quickly surface cost and usage metrics associated with your AWS 
account and visualize them in a human-readable format like a table, csv file, 
or chart. It was created so I could quickly explore and reason about service costs without switching context from the command line.
It's not designed as a replacement for the official AWS COST Explorer CLI but does provide some nice features for visualization and sorting. 


Installation
------------
<hr>

Build from source or download the latest release from the [releases page](https://github.com/cduggn/ccExplorer/releases).

#### From `Homebrew`

```console
$ brew tap cduggn/cduggn

$ brew install ccExplorer
```

#### From `source`:

```console
$ git clone https://github.com/cduggn/ccExplorer.git

$ cd ccExplorer 

$ go run .\cmd\ccexplorer\ccexplorer.go get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -f SERVICE="Amazon DynamoDB"  -l -p csv
```

#### From`docker`:

```console
# download
$ docker pull ghcr.io/cduggn/ccexplorer:v0.3.10

# Container requires AWS Access key, secret, and region
$ docker run -it \
  -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
  -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> \
  -e AWS_REGION=<AWS-REGION> \
  --mount type=bind,source="$(pwd)"/output/,target=/app/output \
  ghcr.io/cduggn/ccexplorer:v0.3.10 get aws -g DIMENSION=OPERATION,
  DIMENSION=SERVICE \
  -l -p chart
  
```

Quick Start
-----------
<hr>

Once `ccExplorer` is installed you can run the help command to see the 
available commands.

```console
$ ccexplorer --help
```
When you invoke a command, `ccExplorer` will use the AWS 
credential chain to authenticate with AWS.

Use the `run-query` command to view and execute a list of preset commands when getting started.

```console
$ ccexplorer run-query
```

For more advanced usage, you can use the `get` command to query AWS Cost and Usage Reports.

```console

#### Sample queries

```console

# Costs grouped by LINKED_ACCOUNT 
$ ccexplorer get aws -g DIMENSION=LINKED_ACCOUNT

# Costs grouped by CommittedThroughput operation and SERVICE
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-10-10 -f OPERATION="CommittedThroughput" -l

# Costs grouped by CommittedThroughput and LINKED_ACCOUNT
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=LINKED_ACCOUNT  -s 2022-10-10 -f OPERATION="CommittedThroughput" -l

# DynamodDB costs grouped by OPERATION
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-10-10 -f SERVICE="Amazon DynamoDB" -l

# All service costs grouped by SERVICE
$ ccexplorer get aws -g DIMENSION=SERVICE -s 2022-10-10

# All service costs grouped by SERVICE and OPERATION
$ ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2022-10-01 -l

# All service costs grouped by SERVICE and OPERATION and sorted in descending order by date
$ ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2023-01-01 -e 2023-02-10 -l -d -m DAILY

# S3 costs grouped by OPERATION 
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-04-04  -f SERVICE="Amazon Simple Storage Service" -l

# Costs grpuped by ApplicationName Cost Allocation Tag
$ ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -l

# Costs grpuped by ApplicationName Cost Allocation Tag and filtered by specific name
$ ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -f TAG="my-project" -l

# S3 costs grouped by SERVICE dimension and ApplicationName Cost Allocation Tag
$ ccexplorer get aws -g DIMENSION=SERVICE,TAG=ApplicationName -f SERVICE="Amazon Simple Storage Service"  -l

# S3 costs grouped by SERVICE dimension and ApplicationName Cost Allocation Tag and filtered by specific name
$ ccexplorer get aws -g DIMENSION=SERVICE,TAG=ApplicationName -f SERVICE="Amazon Simple Storage Service"  -l -f TAG="my-application"

# S3 costs grouped by SERVICE dimension and BucketName Cost Allocation Tag
$ ccexplorer get aws -g DIMENSION=SERVICE,TAG=BucketName -f SERVICE="Amazon Simple Storage Service" -l

# S3 costs grouped by SERVICE dimension and BucketName Cost Allocation Tag and filterred by specific name
$ ccexplorer get aws -g DIMENSION=OPERATION,TAG=BucketName -f SERVICE="Amazon Simple Storage Service" -l -f TAG="my-bucket"

# Costs groupedby OPERATION dimension and ApplicationName Cost Allocation Tag and filtered by PutObject operation
$ ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -f OPERATION="PutObject" -l

# Costs grouped by GetCostAndUsage operation and LINKED_ACCOUNT dimension
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=LINKED_ACCOUNT -s 2022-12-10 -f OPERATION="GetCostAndUsage" -l

# Costs grouped by HOUR and by SERVICE and OPERATION DIMENSIONS
$ ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -l -e 2023-01-27T15:04:05Z -s 2023-01-26T15:04:05Z -m HOURLY

# Costs grouped by DAY and by SERVICE and OPERATION DIMEBSIONS
$ ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -l -e 2023-01-27 -s 2023-01-26 -m DAILY

# Costs exported in CSV format
$ ccexplorer get aws -g DIMENSION=LINKED_ACCOUNT,DIMENSION=OPERATION -l -m DAILY -p csv

# Costs exported to stdout
$ ccexplorer get aws -g DIMENSION=LINKED_ACCOUNT,DIMENSION=OPERATION -l -m DAILY -p stdout

# Costs grouped by MONTH by SERVICE and OPERATION and printed to chart
$ ccexplorer get aws -g DIMENSION=SERVICE, DIMENSION=OPERATION -l -e 2023-01-27 -s 2023-01-26 -m MONTHLY -p chart

# Costs grouped by MONTH by OPERATION and USAGE_TYPE and printed to chart
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=USAGE_TYPE -l -e 2023-01-27 -s 2023-01-26 -m MONTHLY -p chart


```

#### Default settings
If no cost allocation tags have been defined, the  `ccExplorer` can still be 
used to filter and group resources based on their 
AWS resource types. This can be achieved by using the group by and filter 
flags 

- If no billing period is specified, the current calendar month will be used. 
- UnblendedCost is the default cost metric. Other metrics can be specified 
  using the `-i` flag.
- `ccExplorer` prints to stdout by default. The `-p` flag can be used to 
  specify the output format (csv, chart, stdout).
- Results are sorted by default by cost in descending order. The `-d` flag 
  can be used to specify date sorting in descending order.
- Refunds, discounts and credits are applied automatically. The `-l` flag 
  should be used to exclude this behavior.
- When filtering by cost allocation tags (`-f TAG="my-tag"`) a tag must also 
  be specified in the group by flag (`-g TAG=ApplicationName`). This 
  instructs the `ccExplorer` to filter by `ApplicationName=my-tag` .
- Hourly results can be returned by using the `-s` and `-e` flags and 
  providing an ISO 8601 formatted date and time for example `-s 
  2022-10-10T00:00:00Z -e 2022-10-10T23:59:59Z`. 
  

## Additional Information
<hr>

- Cost Explorer accesses data for the last 12 months.
- Cost Explorer charges per paginated request.
- The AWS SDK uses the default credentials provider chain.
- Credits and refunds are automatically applied to Cost Explorer results.
- Cost Explorer API calls can be tracked using CloudTrail. 
- Requests are issued against the `us-east-1` region.
