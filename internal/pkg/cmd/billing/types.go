package billing

type AWSService int64

const (
	AWSServiceEC2 AWSService = iota
	AWSServiceECS
	AWSServiceEKS
	AWSServiceElasticBeanstalk
	AWSServiceElasticLoadBalancing
	AWSServiceElasticMapReduce
	AWSServiceElasticache
	AWSServiceGlacier
	AWSServiceKinesis
	AWSServiceLambda
	AWSServiceRDS
	AWSServiceRedshift
	AWSServiceS3
	AWSServiceSES
	AWSServiceSNS
	AWSServiceSQS
	AWSServiceStorageGateway
	AWSServiceVPC
	AWSAmplify
	AWSAppSync
	AWSAppStream
	AWSAConfig
	AWSDataTransfer
	AWSKeyManagementService
	AWSLambda
	AWSSecretsManager
	AWSStepFunctions
	AWSAPIGateway
	AWSCloudFront
	AWSCognito
	AWSDynamoDB
	AWSElasticFileSystem
	AWSRegistrar
	AWSRoute53
	Refund
	Tax
)

type billable interface {
	print(Service)
	billingDetails(string) Service
	total() float64
}

type Billable struct {
	Services map[string]Service
	Start    string
	End      string
}

type Service struct {
	Name    string
	Metrics []Metrics
}

type Metrics struct {
	Name   string
	Amount string
	Unit   string
}
