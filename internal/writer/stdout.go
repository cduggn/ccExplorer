package writer

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

var (
	tableDivider = table.Row{"", "", "",
		"", "", "", "",
		"",
		"", ""}
	costAndUsageHeader = table.Row{"Rank", "Dimension/Tag", "Dimension/Tag",
		"Metric Name", "Amount", "Rounded",
		"Unit",
		"Granularity",
		"Start",
		"End"}
	costAndUsageTableFooter = func(t string) table.Row {
		return table.
			Row{"", "",
			"",
			"",
			"Cost",
			t, "", "", "", ""}
	}
	forecastedHeader = table.Row{"Start", "End", "Mean Value",
		"Prediction Interval LowerBound",
		"Prediction Interval UpperBound", "Unit", "Total"}
	forecastedTableFooter = func(filter string, unit string,
		amount string) table.Row {
		return table.Row{"FilteredBy", filter, "", "", "",
			unit,
			amount}
	}
	tableStyleColors = table.StyleColoredGreenWhiteOnBlack
)

type CostAndUsageTable struct {
	Table table.Writer
}

type ForecastTable struct {
	Table table.Writer
}

func NewStdoutWriter(variant string) (types.Table, error) {
	switch variant {
	case "forecast":
		return ForecastTable{
			Table: table.NewWriter(),
		}, nil
	case "costAndUsage":
		return CostAndUsageTable{
			Table: table.NewWriter(),
		}, nil
	}

	return nil, fmt.Errorf("unknown table type: %s", variant)
}

func (c CostAndUsageTable) Writer(output interface{}) {
	outputType := output.(types.CostAndUsageStdoutType)
	c.Style()
	c.Header()
	rows := CostUsageToRows(outputType.Services, outputType.Granularity)

	c.AddRows(rows.Rows)
	c.Table.AppendRow(tableDivider)
	c.Footer(costAndUsageTableFooter(rows.Total))
	c.Table.Render()
}

func (c CostAndUsageTable) Style() {
	c.Table.SetOutputMirror(os.Stdout)
	c.Table.SetColumnConfigs(
		[]table.ColumnConfig{
			{Number: 6, WidthMax: 8},
		})
	c.Table.SetStyle(tableStyleColors)

	c.Table.SuppressEmptyColumns()
}

func (c CostAndUsageTable) Header() {
	c.Table.AppendHeader(costAndUsageHeader)
}

func (c CostAndUsageTable) AddRows(rows []table.Row) {
	c.Table.AppendRows(rows)
}

func (c CostAndUsageTable) Footer(row table.Row) {
	c.Table.AppendFooter(row)
}

func (f ForecastTable) Writer(output interface{}) {
	outputType := output.(types.ForecastStdoutType)
	//f.Table.SuppressEmptyColumns()
	f.Style()
	f.Header()
	rows := ForecastToRows(outputType)
	f.AddRows(rows)

	f.Footer(forecastedTableFooter(outputType.FilteredBy,
		outputType.Total.Unit, outputType.Total.Amount))

	f.Table.Render()
}

func (f ForecastTable) Style() {
	f.Table.SetOutputMirror(os.Stdout)
	f.Table.SetStyle(tableStyleColors)
}

func (f ForecastTable) Footer(row table.Row) {
	f.Table.AppendFooter(row)
}

func (f ForecastTable) AddRows(rows []table.Row) {
	f.Table.AppendRows(rows)
}

func (f ForecastTable) Header() {
	f.Table.AppendHeader(forecastedHeader)
}

func CostUsageToRows(s []types.Service, granularity string) types.CostAndUsage {
	var rows []table.Row
	var total float64
	for index, v := range s {
		for _, m := range v.Metrics {

			if index%10 == 0 {
				rows = append(rows, tableDivider)
			}
			if m.Unit == "USD" {
				total += m.NumericAmount
			}

			tempRow := table.Row{index + 1, v.Keys[0],
				utils.ReturnIfPresent(v.Keys),
				m.Name, m.Amount, fmt.Sprintf("%.2f",
					m.NumericAmount),
				m.Unit,
				granularity,
				v.Start, v.End}

			rows = append(rows, tempRow)
		}
	}
	totalFormatted := fmt.Sprintf("$%.2f", total)
	return types.CostAndUsage{Rows: rows, Total: totalFormatted}
}

func ForecastToRows(r types.ForecastStdoutType) []table.Row {
	var rows []table.Row
	for _, v := range r.Forecast {
		tempRow := table.Row{v.TimePeriod.Start,
			v.TimePeriod.End, v.MeanValue, v.PredictionIntervalUpperBound,
			v.PredictionIntervalLowerBound}

		rows = append(rows, tempRow)
	}
	return rows
}
