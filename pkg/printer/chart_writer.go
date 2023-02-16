package printer

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

func (Renderer) Charts(r CostAndUsageOutputType) error {
	page := components.NewPage()
	page.PageTitle = "Cost and Usage Report"

	p := buildPieCharts(r)
	for _, chart := range p {
		page.AddCharts(chart)
	}

	f, err := createOutputFile()
	if err != nil {
		return PrinterError{
			msg: "Failed creating ccexplorer.html: " + err.Error(),
		}
	}

	err = renderCharts(page, f)
	if err != nil {
		return PrinterError{
			msg: "Failed rendering chart: " + err.Error(),
		}
	}
	return nil
}

func buildPieCharts(r CostAndUsageOutputType) []*charts.Pie {

	var pieC []*charts.Pie
	dimensions := r.Dimensions
	if len(dimensions) > 1 {
		for index, dimension := range dimensions {
			pieC = append(pieC, DefinePieChartProperties(r.Services, dimension, index,
				r.Granularity, r.Start, r.End))
		}
	} else {
		pieC = append(pieC, DefinePieChartProperties(r.Services, dimensions[0], 0, r.Granularity, r.Start, r.End))
	}
	return pieC

}

func DefinePieChartProperties(s map[int]Service, d string, index int,
	granularity string, start string, end string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: CreateTitle(d),
			TitleStyle: &opts.TextStyle{
				FontSize:  20,
				FontStyle: "normal",
				Padding:   80,
			},
			Bottom: "0",
			Subtitle: fmt.Sprintf("Granularity: %s Start: %s, End: %s",
				granularity, start,
				end),
		}),
	)

	pie.AddSeries("pie", PopulatePieDate(s, index)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{b} : {c}",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"30%", "60%"},
			}),
		)

	pie.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{
			Padding: 10,
		}),
	)

	return pie
}

func createOutputFile() (*os.File, error) {
	return os.Create("./output/ccexplorer.html")
}

func renderCharts(p *components.Page, f *os.File) error {
	return p.Render(io.MultiWriter(f))
}
