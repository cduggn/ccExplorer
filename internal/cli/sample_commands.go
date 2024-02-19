package cli_new

import (
	"github.com/cduggn/ccexplorer/internal/types"
)

const (
	CostAndUsageExamples = `
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
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2022-10-01 -l

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

  # Costs grouped by MONTH by SERVICE and OPERATION and printed to chart
  ccexplorer get aws -g DIMENSION=SERVICE, DIMENSION=OPERATION -l -e 2023-01-27 -s 2023-01-26 -m MONTHLY -p chart
 
  # Costs grouped by MONTH by OPERATION and USAGE_TYPE and printed to chart
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=USAGE_TYPE -l -e 2023-01-27 -s 2023-01-26 -m MONTHLY -p chart

  # All service costs grouped by SERVICE and OPERATION and sorted in descending order by date
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2023-01-01 -e 2023-02-10 -l -d

`
	ForecastExamples = `
  # Service forecast for the next 30 days
  ccexplorer get aws forecast -p 95 -g MONTHLY

  # S3 cost forecast for the next 30 days
  ccexplorer get aws forecast -f SERVICE="Amazon Simple Storage Service"  -p 95 -g MONTHLY
  
  # DynamoDB cost forecast for PutObject operations for the next 30 days
  ccexplorer get aws forecast -f SERVICE="Amazon DynamoDB",OPERATION="CommittedThroughput"  -p 95 -g MONTHLY
  
`
)

func PresetList() []types.PresetParams {
	p := []types.PresetParams{
		{
			Alias:             "Costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "",
			Filter:            map[string]string{},
			FilterByTag:       false,
			FilterByDimension: false,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=LINKED_ACCOUNT)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=LINKED_ACCOUNT -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "Costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "",
			Filter:            map[string]string{},
			FilterByTag:       false,
			FilterByDimension: false,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=USAGE_TYPE)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=USAGE_TYPE -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "Costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "",
			Filter:            map[string]string{},
			FilterByTag:       false,
			FilterByDimension: false,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=OPERATION)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "S3 costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=OPERATION)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -f SERVICE=\"Amazon Simple Storage" +
				" Service\"]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "S3 costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=USAGE_TYPE)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=USAGE_TYPE -f SERVICE=\"Amazon Simple Storage" +
				" Service\" -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "S3 costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=LINKED_ACCOUNT)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=LINKED_ACCOUNT -f SERVICE=\"Amazon Simple Storage" +
				" Service\" -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "DynamoDB costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon DynamoDB"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=OPERATION)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -f SERVICE=\"Amazon DynamoDB\" -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "DynamoDB costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon DynamoDB"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=USAGE_TYPE)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=USAGE_TYPE -f SERVICE=\"Amazon DynamoDB\" -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "DynamoDB costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon DynamoDB"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=LINKED_ACCOUNT)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=LINKED_ACCOUNT -f SERVICE=\"Amazon DynamoDB\" -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
		{
			Alias:             "Costs grouped by GetCostAndUsage OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "Name",
			Filter:            map[string]string{"OPERATION": "GetCostAndUsage"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			Description: []string{"(Dimension=SERVICE",
				"Dimension=LINKED_ACCOUNT)"},
			CommandSyntax: "[ ccexplorer get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -f OPERATION=\"GetCostAndUsage\" -l ]",
			Granularity: "MONTHLY",
			PrintFormat: "stdout",
			Metric:      []string{"UnblendedCost"},
		},
	}
	return p
}
