package writer

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	costexplorertypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnomaliesToTableTransformer(t *testing.T) {
	input := types.AnomaliesPrintData{
		Anomalies: &costexplorer.GetAnomaliesOutput{
			Anomalies: []costexplorertypes.Anomaly{
				{
					AnomalyId:        aws.String("anomaly-1"),
					MonitorArn:       aws.String("arn:aws:ce::123456789012:anomalymonitor/monitor-id"),
					DimensionValue:   aws.String("Amazon EC2"),
					AnomalyStartDate: aws.String("2024-01-01"),
					AnomalyEndDate:   aws.String("2024-01-03"),
					AnomalyScore:     &costexplorertypes.AnomalyScore{CurrentScore: 87.5, MaxScore: 90},
					Impact:           &costexplorertypes.Impact{MaxImpact: 123.45},
					Feedback:         costexplorertypes.AnomalyFeedbackTypeYes,
				},
				{
					// Optional fields left unset, as a real anomaly can have.
					AnomalyId:    aws.String("anomaly-2"),
					MonitorArn:   aws.String("arn:aws:ce::123456789012:anomalymonitor/monitor-id"),
					AnomalyScore: &costexplorertypes.AnomalyScore{CurrentScore: 10, MaxScore: 10},
					Impact:       &costexplorertypes.Impact{MaxImpact: 5},
				},
			},
		},
		Filters: []string{"MonitorArn: arn:aws:ce::123456789012:anomalymonitor/monitor-id"},
	}

	out, err := NewAnomaliesToTableTransformer().Transform(input)
	require.NoError(t, err)

	assert.Equal(t, 2, out.RowCount)
	require.Len(t, out.Rows, 2)

	first := out.Rows[0]
	assert.Equal(t, "anomaly-1", first[0])
	assert.Equal(t, "Amazon EC2", first[2])
	assert.Equal(t, "2024-01-01", first[3])
	assert.Equal(t, "$123.45", first[6])
	assert.Equal(t, "YES", first[7])

	// The transformer must not panic when optional pointer fields are nil.
	second := out.Rows[1]
	assert.Equal(t, "", second[2])
	assert.Equal(t, "", second[3])
	assert.Equal(t, "", second[7])

	assert.Equal(t, "MonitorArn: arn:aws:ce::123456789012:anomalymonitor/monitor-id", out.FilterInfo)
}

func TestAnomaliesToTableTransformerEmptyResult(t *testing.T) {
	out, err := NewAnomaliesToTableTransformer().Transform(types.AnomaliesPrintData{
		Anomalies: &costexplorer.GetAnomaliesOutput{},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, out.RowCount)
	assert.Empty(t, out.FilterInfo)
}
