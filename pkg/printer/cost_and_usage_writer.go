package printer

func CostAndUsageToStdout(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) {
	sortedServices := sortFn(r.Services)

	t := CreateTable(costAndUsageHeader)

	granularity := r.Granularity

	rows := CostUsageToRows(sortedServices, granularity)

	t.AppendRows(rows.Rows)
	t.AppendRow(tableDivider)
	t.AppendRow(costAndUsageTableFooter(rows.Total))

	t.Render()
}

func CostAndUsageToCSV(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	f, err := NewCSVFile(OutputDir, csvFileName)
	if err != nil {
		return PrinterError{
			msg: "Error creating CSV file: " + err.Error()}
	}
	defer f.Close()

	rows := ToRows(r.Services, r.Granularity)

	err = CSVWriter(f, csvheader, rows)
	if err != nil {
		return nil
	}

	return nil
}

func CostAndUsageToChart(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	builder := ChartBuilder{}
	charts, err := builder.NewCharts(r)
	if err != nil {
		return err
	}

	err = ChartWriter(charts)
	if err != nil {
		return err
	}
	return nil
}

func CostAndUsageToOpenAI(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	rows := ToRows(r.Services, r.Granularity)

	data := BuildPromptText(rows)

	resp, err := SummarizeWIthAI(r.OpenAIAPIKey, data)
	if err != nil {
		return err
	}

	f, err := NewFile(OutputDir, aiFileName)
	if err != nil {
		return PrinterError{
			msg: "Failed creating AI HTML: " + err.Error(),
		}
	}
	defer f.Close()

	err = AIWriter(f, resp.Choices[0].Text)
	if err != nil {
		return err
	}

	return nil
}
