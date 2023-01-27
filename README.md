###### ccExplorer is in Alpha

<h1 align="center"><code>ccExplorer</code></h1>

<hr>
<div align="center">
<a href="https://github.com/cduggn/ccExplorer/actions" 
alt="goreleaser status">
<img src="https://github.com/cduggn/ccExplorer/actions/workflows/release.yml/badge.svg">
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
cost of your cloud resources. It's not 
designed as a replacement for the official AWS CLI and does not offer the 
same exhaustive search option. It does however return results in a more
human-readable format, and orders them by cost in descending order.
It's primary use case is to surface costs based on pre-defined cost allocation tags. 

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

$ go run .\cmd\ccexplorer\ccexplorer.go get aws -g DIMENSION=SERVICE,
DIMENSION=OPERATION -f SERVICE="Amazon DynamoDB"  -l
```

#### From`docker`:

```console
# download

$ docker pull ghcr.io/cduggn/ccexplorer:v0.2.0

# Container requires AWS Access key, secret, and region

$ docker run -it \
  -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
  -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> \
  -e AWS_REGION=<AWS-REGION> \
  ghcr.io/cduggn/ccexplorer:v0.2.0 get aws -g DIMENSION=OPERATION,
  DIMENSION=SERVICE -l 
  
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

Presets are the most common way to use `ccExplorer`. Presets are a 
set of pre-defined cost and usage queries that can be used to quickly get a 
sense of the cost of your cloud resources. 

```console
$ ccexplorer run-query
```

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
$ ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2022-10- -l

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
ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -f OPERATION="PutObject" -l

# Costs grouped by GetCostAndUsage operation and LINKED_ACCOUNT dimension
$ ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=LINKED_ACCOUNT -s 2022-12-10 -f OPERATION="GetCostAndUsage" -l

```

#### Default settings
If no cost allocation tags have been defined, the  `ccExplorer` can still be 
used to filter and group resources based on their 
AWS resource types. This can be achieved by using the group by and filter 
flags 

- If no billing period is specified, thecurrent calendar month will be used. 
- Results are sorted by cost in descending order.
- Refunds, discounts and credits are applied automatically. The `-l` flag 
  should be used to exclude this behavior.
- When filtering by cost allocation tags (`-f TAG="my-tag"`) a tag must also 
  be specified in the group by flag (`-g TAG=ApplicationName`). This 
  instructs the `ccExplorer` to filter by `ApplicationName=my-tag` .
  

## Additional Information
<hr>

- Cost Explorer accesses data for the last 12 months.
- Cost Explorer charges per paginated request.
- The AWS SDK uses the default credentials provider chain.
- Credits and refunds are automatically applied to Cost Explorer results.
- Cost Explorer API calls can be tracked using CloudTrail. 
- Requests are issued against the `us-east-1` region.
