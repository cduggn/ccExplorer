package cost_and_usage

func ValidatePrintFormat(f string, key string) bool {
	return IsValidPrintFormat(f) && hasOpenAIAPIKey(f, key)
}

func IsValidPrintFormat(f string) bool {
	return f == "stdout" || f == "csv" || f == "chart" || f == "gpt3"
}

func hasOpenAIAPIKey(f string, key string) bool {
	if f == "gpt3" && key == "" {
		return false
	}
	return true
}

func IsValidGranularity(g string) bool {
	return g == "DAILY" || g == "MONTHLY" || g == "HOURLY"
}

func IsValidMetric(m string) bool {
	return m == "AmortizedCost" || m == "BlendedCost" || m == "NetAmortizedCost" ||
		m == "NetUnblendedCost" || m == "NormalizedUsageAmount" || m == "UnblendedCost" ||
		m == "UsageQuantity"
}

func SortByFn(sortByDate bool) string {
	if sortByDate {
		return "date"
	}
	return "cost"
}
