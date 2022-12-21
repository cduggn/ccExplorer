package aws

import (
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
