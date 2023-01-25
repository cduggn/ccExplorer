package presets

func PresetList() []PresetParams {
	p := []PresetParams{
		{
			Alias:             "Costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "",
			Filter:            map[string]string{},
			FilterByTag:       false,
			FilterByDimension: false,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=LINKED_ACCOUNT -l",
		},
		{
			Alias:             "Costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "",
			Filter:            map[string]string{},
			FilterByTag:       false,
			FilterByDimension: false,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=USAGE_TYPE -l",
		},
		{
			Alias:             "Costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "",
			Filter:            map[string]string{},
			FilterByTag:       false,
			FilterByDimension: false,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -l",
		},
		{
			Alias:             "S3 costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -f SERVICE=\"Amazon Simple Storage" +
				" Service\"",
		},
		{
			Alias:             "S3 costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=USAGE_TYPE -f SERVICE=\"Amazon Simple Storage" +
				" Service\"",
		},
		{
			Alias:             "S3 costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=LINKED_ACCOUNT -f SERVICE=\"Amazon Simple Storage" +
				" Service\"",
		},
		{
			Alias:             "DynamoDB costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon DynamoDB"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=OPERATION -f SERVICE=\"Amazon DynamoDB\"",
		},
		{
			Alias:             "DynamoDB costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon DynamoDB"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=USAGE_TYPE -f SERVICE=\"Amazon DynamoDB\"",
		},
		{
			Alias:             "DynamoDB costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon DynamoDB"},
			FilterByTag:       false,
			FilterByDimension: true,
			ExcludeDiscounts:  true,
			CommandSyntax: "get aws -g DIMENSION=SERVICE," +
				"DIMENSION=LINKED_ACCOUNT -f SERVICE=\"Amazon DynamoDB\"",
		},
	}
	return p
}
