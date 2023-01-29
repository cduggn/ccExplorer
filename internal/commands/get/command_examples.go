package get

const (
	costAndUsageExamples = `
  # Costs grouped by LINKED_ACCOUNT 
  ccexplorer get aws -g DIMENSION=LINKED_ACCOUNT
  
  # Costs grouped by CommittedThroughput operation and SERVICE
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-10-10 -f OPERATION="CommittedThroughput" -l

  # Costs grouped by CommittedThroughput and LINKED_ACCOUNT
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=LINKED_ACCOUNT  -s 2022-10-10 -f OPERATION="CommittedThroughput" -l

  # DynamodDB costs grouped by OPERATION
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-10-10 -f SERVICE="Amazon DynamoDB" -l

  # All service costs grouped by SERVICE
  ccexplorer get aws -g DIMENSION=SERVICE -s 2022-10-10

  # All service costs grouped by SERVICE and OPERATION
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2022-10- -l

  # S3 costs grouped by OPERATION
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-04-04  -f SERVICE="Amazon Simple Storage Service" -l

  # Costs grpuped by ApplicationName Cost Allocation Tag
  ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -l
 
  # Costs grouped by HOUR by SERVICE and OPERATION
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -l -e 2023-01-27T15:04:05Z -s 2023-01-26T15:04:05Z -m HOURLY

  # Costs grouped by DAY by SERVICE and OPERATION
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -l -e 2023-01-27 -s 2023-01-26 -m DAILY

  # Costs grouped by DAY by SERVICE and OPERATION and printed to CSV
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -l -e 2023-01-27 -s 2023-01-26 -m DAILY -p csv
`
	forecastExamples = `
  # Service forecast for the next 30 days
  ccexplorer get aws forecast -p 95 -g MONTHLY

  # S3 cost forecast for the next 30 days
  ccexplorer get aws forecast -f SERVICE="Amazon Simple Storage Service"  -p 95 -g MONTHLY
  
  # DynamoDB cost forecast for PutObject operations for the next 30 days
  ccexplorer get aws forecast -f SERVICE="Amazon DynamoDB",OPERATION="CommittedThroughput"  -p 95 -g MONTHLY
  
`
)
