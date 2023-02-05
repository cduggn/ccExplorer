package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func pieRadius(r CostAndUsageOutputType) []*charts.Pie {

	var pies []*charts.Pie
	dimensions := r.Dimensions
	if len(dimensions) > 1 {
		for index, dimension := range dimensions {
			pies = append(pies, RenderPieChart(r.Services, dimension, index,
				r.Granularity, r.Start, r.End))
		}
	} else {
		pies = append(pies, RenderPieChart(r.Services, dimensions[0], 0, r.Granularity, r.Start, r.End))
	}
	return pies

}

func RenderPieChart(s map[int]Service, d string, index int,
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

func (Renderer) Charts(r CostAndUsageOutputType) error {
	page := components.NewPage()

	charts := pieRadius(r)

	for _, chart := range charts {
		page.AddCharts(chart)
	}
	page.PageTitle = "Cost and Usage Report"

	f, err := os.Create("./output/ccexplorer.html")
	if err != nil {
		return PrinterError{
			msg: "Failed creating ccexplorer.html: " + err.Error(),
		}
	}
	err = page.Render(io.MultiWriter(f))
	if err != nil {
		return PrinterError{
			msg: "Failed rendering chart: " + err.Error(),
		}
	}
	return nil
}
