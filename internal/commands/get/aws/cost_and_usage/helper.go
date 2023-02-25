package cost_and_usage

func SortByFn(sortByDate bool) string {
	if sortByDate {
		return "date"
	}
	return "cost"
}
