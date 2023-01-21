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
ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -l
```

This will return a list of all AWS services and operations that have been
used in the specified billing period. If no billing period is specified, the
current calendar month will be used. Results are sorted by cost in 
descending order and refunds, discounts and credits are excluded.

```sh
ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -f TAG="my-project"
```
This will return a list of all AWS operations that have been used in the 
specified billing period for the specified project. The `-f` flag can be
used to filter results based on the value of a cost allocation tag. If no 
filter value is specified, all resources will be returned. Results are  
sorted by cost in descending order.

```sh
ccexplorer get aws -g DIMENSION=SERVICE,TAG=ApplicationName -f SERVICE="Amazon Simple Storage Service"  -l
```

This will return a list of costs for S3 buckets that have been used in the
specified billing period and that have been tagged with the ApplicationName
tag. Results are sorted by cost in descending order.

```sh
ccexplorer get aws -g DIMENSION=SERVICE,TAG=ApplicationName -f SERVICE="Amazon Simple Storage Service"  -l -f TAG="my-application"
```

This will return a list of costs for the specified application that have
been used in the specified billing period. Results are sorted by cost in
descending order.

```sh
ccexplorer get aws -g DIMENSION=SERVICE,TAG=BucketName -f SERVICE="Amazon Simple Storage Service"
```

This will return a list of costs for S3 buckets filtered by the bucket name
tag. Results are sorted by cost in descending order.


Installation
------------

Build from source or download the latest release from the [releases page](https://github.com/cduggn/ccExplorer/releases). 

### Run 

#### From Homebrew

```sh
brew tap cduggn/ccExplorer
brew install ccExplorer
```

#### From source: 

```sh
git clone https://github.com/cduggn/ccExplorer.git
cd ccExplorer 
go run .\cmd\ccexplorer\ccexplorer.go get aws -d SERVICE -d OPERATION -u SERVICE="Amazon DynamoDB"  -c
```

#### From`docker`:

```sh
# download

docker pull cdugga/ccexplorer:0.1.0

# run requires AWS Access key and region to be set

docker run \
  -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
  -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> \
  -e AWS_REGION=<AWS-REGION> \
  cdugga/ccexplorer:0.1.0 get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -l 
```

## Considerations when using Cost Explorer

- Cost Explorer accesses data for the last 12 months.
- Cost Explorer charges per paginated request.
- The AWS SDK uses the default credentials provider chain.
- Credits and refunds are automatically applied to Cost Explorer results.`.
- Cost Explorer API calls can be tracked using CloudTrail. Requests are issued against us-east-1.
