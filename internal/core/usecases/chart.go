package usecases

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/util"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
)

var (
	chartFileName = "ccexplorer_chart.html"
)

func WriteToChart(p *components.Page) error {

	f, err := util.NewFile(OutputDir, chartFileName)
	if err != nil {
		return model.Error{
			Msg: "Failed creating chart HTML file: " + err.Error(),
		}
	}
	defer f.Close()
	return p.Render(io.MultiWriter(f))
}

func (Builder) NewCharts(r model.InputType) (*components.Page,
	error) {
	page := components.NewPage()
	page.PageTitle = "Cost and Usage Report"

	p := buildPieCharts(r)
	for _, chart := range p {
		page.AddCharts(chart)
	}

	return page, nil
}

func buildPieCharts(r model.InputType) []*charts.Pie {

	var pieC []*charts.Pie
	dimensions := r.Dimensions
	if len(dimensions) > 1 {
		for index, dimension := range dimensions {
			pieC = append(pieC, definePieChartProperties(r.Services, dimension, index,
				r.Granularity, r.Start, r.End))
		}
	} else {
		pieC = append(pieC, definePieChartProperties(r.Services, dimensions[0], 0, r.Granularity, r.Start, r.End))
	}
	return pieC

}

func definePieChartProperties(s []model.Service, d string, index int,
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

func PopulatePieDate(services []model.Service, key int) []opts.
	PieData {
	items := make([]opts.PieData, 0)

	//services := SortServicesByMetricAmount(s)

	for index, v := range services {
		if index < 15 {
			items = append(items, opts.PieData{Name: v.Keys[key],
				Value: v.Metrics[0].NumericAmount})
		}

	}
	return items
}

func CreateTitle(dimension string) string {
	return fmt.Sprintf("Pie chart for dimension: [ %s ]", dimension)
}
