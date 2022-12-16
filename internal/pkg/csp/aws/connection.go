package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

var (
	groupByDimension = func(dimensions []string) []types.GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		return groups
	}
	groupByTag = func(tag string) []types.GroupDefinition {
		return []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeTag,
				Key:  aws.String(tag),
			},
		}
	}
	groupByTagAndDimension = func(tag string, dimensions []string) []types.GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		groups = append(groups, types.GroupDefinition{
			Type: types.GroupDefinitionTypeTag,
			Key:  aws.String(tag),
		})
		return groups
	}
	filterCredits = func() *types.Expression {
		return &types.Expression{
			Not: &types.Expression{
				Dimensions: &types.DimensionValues{
					Key:    "RECORD_TYPE",
					Values: []string{"Refund", "Credit"},
				},
			},
		}
	}
	filterByTag = func(tag string, value string) *types.Expression {
		return &types.Expression{
			Tags: &types.TagValues{
				Key:    aws.String(tag),
				Values: []string{value},
			},
		}
	}
	conn *storage.CostDataStorage
)

func init() {
	newConnection()
}

func newConnection() {
	conn = &storage.CostDataStorage{}
	err := conn.New("./cloudcost.db")
	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *CostAndUsageReport) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dimension/Tag", "Dimension/Tag", "Metric Name", "Amount", "Unit", "Granularity", "Start", "End"})
	var total float64
	for _, m := range c.Services {
		for _, v := range m.Metrics {
			if v.Unit == "USD" {
				total += ConvertToFloat(v.Amount)
			}
			tempRow := table.Row{m.Keys[0], isEmpty(m.Keys), v.Name, v.Amount, v.Unit, c.Granularity, m.Start, m.End}
			t.AppendRow(tempRow)

			conn.Insert(storage.CostDataInsert{
				Dimension:   m.Keys[0],
				Dimension2:  "",
				Tag:         "",
				MetricName:  "",
				Amount:      v.Amount,
				Unit:        v.Unit,
				Granularity: c.Granularity,
				StartDate:   m.Start,
				EndDate:     m.End,
			})
		}
	}
	totalHeaderRow := table.Row{"", "", "", "", "", "", "", ""}
	totalRow := table.Row{"", "", "TOTAL COST", total, "", "", "", ""}
	t.AppendRow(totalHeaderRow)
	t.AppendRow(totalRow)
	t.Render()
}

func isEmpty(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}

}
func filter(req CostAndUsageRequestType) *types.Expression {
	expression := &types.Expression{}

	if req.ExcludeCredits && req.IsFilterEnabled {
		expression.And = []types.Expression{*filterCredits(), *filterByTag(req.Tag, req.TagFilterValue)}
	} else if req.ExcludeCredits {
		expression = filterCredits()
	} else if req.IsFilterEnabled {
		expression = filterByTag(req.Tag, req.TagFilterValue)
	} else {
		return nil
	}
	return expression
}

func groupBy(req CostAndUsageRequestType) []types.GroupDefinition {
	if req.Tag != "" && len(req.GroupBy) == 1 {
		return groupByTagAndDimension(req.Tag, req.GroupBy)
	} else if req.Tag != "" {
		return groupByTag(req.Tag)
	} else {
		return groupByDimension(req.GroupBy)
	}

}

func (c *CostAndUsageReport) CurateReport(output *costexplorer.GetCostAndUsageOutput) {
	count := 0
	for _, v := range output.ResultsByTime {
		c.Start = *v.TimePeriod.Start
		c.End = *v.TimePeriod.End
		for _, g := range v.Groups {
			keys := make([]string, 0)
			service := Service{
				Start: c.Start,
				End:   c.End,
			}
			keys = append(keys, g.Keys...)

			for key, m := range g.Metrics {
				metrics := Metrics{
					Name:   key,
					Amount: *m.Amount,
					Unit:   *m.Unit,
				}
				service.Metrics = append(service.Metrics, metrics)
			}
			service.Keys = keys
			c.Services[count] = service
			count++
		}

	}
}
