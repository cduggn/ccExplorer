package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToSlice(t *testing.T) {
	mockDimensionValuesOutput := costexplorer.GetDimensionValuesOutput{
		DimensionValues: []types.DimensionValuesWithAttributes{
			{
				Attributes: map[string]string{
					"ServiceName": "Amazon Elastic Compute Cloud - Compute",
				},
				Value: aws.String("value"),
			},
		},
	}
	slice := ToSlice(mockDimensionValuesOutput)
	assert.Equal(t, "value", slice[0])
	assert.Equal(t, slice, []string{"value"})
}

type mockGetDimensionValuesAPI func(ctx context.Context,
	params *costexplorer.GetDimensionValuesInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetDimensionValuesOutput, error)

func (m mockGetDimensionValuesAPI) GetDimensionValues(ctx context.Context, params *costexplorer.GetDimensionValuesInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetDimensionValuesOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetDimensionValues(t *testing.T) {
	mockDimensionValuesOutput := costexplorer.GetDimensionValuesOutput{
		DimensionValues: []types.DimensionValuesWithAttributes{
			{
				Attributes: map[string]string{
					"ServiceName": "Amazon Elastic Compute Cloud - Compute",
				},
				Value: aws.String("value"),
			},
		},
	}
	mockAPI := mockGetDimensionValuesAPI(func(ctx context.Context, params *costexplorer.GetDimensionValuesInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetDimensionValuesOutput, error) {
		return &mockDimensionValuesOutput, nil
	})
	api := &APIClient{}
	dimensionValues, err := api.GetDimensionValues(context.
		TODO(), mockAPI, GetDimensionValuesRequest{
		Dimension: "SERVICE",
		Time:      Time{Start: "2020-01-01", End: "2020-01-01"},
	})
	assert.NoError(t, err)
	assert.Equal(t, dimensionValues, []string{"value"})
}
