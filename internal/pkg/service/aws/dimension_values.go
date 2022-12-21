package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type GetDimensionValuesAPI interface {
	GetDimensionValues(ctx context.Context,
		params *costexplorer.GetDimensionValuesInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetDimensionValuesOutput, error)
}

func (*APIClient) GetDimensionValues(ctx context.Context,
	api GetDimensionValuesAPI, d GetDimensionValuesRequest) (
	[]string, error) {

	dimensionValues, err := api.GetDimensionValues(context.TODO(),
		&costexplorer.GetDimensionValuesInput{
			Dimension: types.Dimension(d.Dimension),
			TimePeriod: &types.DateInterval{
				Start: aws.String(d.Time.Start),
				End:   aws.String(d.Time.End),
			},
		})

	if err != nil {
		return nil, APIError{
			msg: "Error while fetching Dimension Values for Dimension from AWS",
		}
	}

	ds := ToSlice(*dimensionValues)

	return ds, nil
}

func ToSlice(d costexplorer.GetDimensionValuesOutput) []string {
	var servicesSlice []string
	for _, service := range d.DimensionValues {
		servicesSlice = append(servicesSlice, *service.Value)
	}
	return servicesSlice
}
