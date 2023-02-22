package chart

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

var (
	chartFileName = "ccexplorer_chart.html"
	OutputDir     = "./output"
)

func Writer(p *components.Page) error {

	f, err := newFile(OutputDir, chartFileName)
	if err != nil {
		return Error{
			msg: "Failed creating chart HTML file: " + err.Error(),
		}
	}
	defer f.Close()
	return p.Render(io.MultiWriter(f))
}

func (Builder) NewCharts(r InputType) (*components.Page,
	error) {
	page := components.NewPage()
	page.PageTitle = "Cost and Usage Report"

	p := buildPieCharts(r)
	for _, chart := range p {
		page.AddCharts(chart)
	}

	return page, nil
}

func buildPieCharts(r InputType) []*charts.Pie {

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

func definePieChartProperties(s []Service, d string, index int,
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

// todo remove duplication
func newFile(dir string, file string) (*os.File, error) {
	filePath := buildOutputFilePath(dir, file)
	return os.Create(filePath)
}

// todo remove duplication
func buildOutputFilePath(dir string, fileName string) string {
	return dir + "/" + fileName
}

func PopulatePieDate(services []Service, key int) []opts.
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