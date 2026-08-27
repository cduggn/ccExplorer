package awsservice

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	types2 "github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeCostExplorer returns a canned page per call and records what it was
// asked for.
type fakeCostExplorer struct {
	pages     []*costexplorer.GetCostAndUsageOutput
	err       error
	calls     int
	seenToken []string
	seenCtx   []context.Context

	anomaliesPages     []*costexplorer.GetAnomaliesOutput
	anomaliesErr       error
	anomaliesCalls     int
	anomaliesSeenToken []string
	// anomaliesSeenParams snapshots the fields tests care about at call time,
	// not the *GetAnomaliesInput pointer itself: production code mutates and
	// reuses a single input struct across pages, so storing the pointer would
	// make every entry alias the same final state after the loop ends.
	anomaliesSeenParams []anomaliesSeenParam
}

type anomaliesSeenParam struct {
	monitorArn string
	feedback   types.AnomalyFeedbackType
	maxResults int32
}

func (f *fakeCostExplorer) GetCostAndUsage(ctx context.Context,
	params *costexplorer.GetCostAndUsageInput,
	_ ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {

	f.calls++
	f.seenCtx = append(f.seenCtx, ctx)
	if params.NextPageToken == nil {
		f.seenToken = append(f.seenToken, "")
	} else {
		f.seenToken = append(f.seenToken, *params.NextPageToken)
	}

	if f.err != nil {
		return nil, f.err
	}
	if f.calls <= len(f.pages) {
		return f.pages[f.calls-1], nil
	}
	// Beyond the canned pages, keep handing out a token forever.
	return &costexplorer.GetCostAndUsageOutput{
		NextPageToken: aws.String("endless"),
	}, nil
}

func (f *fakeCostExplorer) GetCostForecast(_ context.Context,
	_ *costexplorer.GetCostForecastInput,
	_ ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error) {
	return nil, errors.New("not used")
}

func (f *fakeCostExplorer) GetAnomalies(_ context.Context,
	params *costexplorer.GetAnomaliesInput,
	_ ...func(*costexplorer.Options)) (*costexplorer.GetAnomaliesOutput, error) {

	f.anomaliesCalls++
	seen := anomaliesSeenParam{feedback: params.Feedback}
	if params.MonitorArn != nil {
		seen.monitorArn = *params.MonitorArn
	}
	if params.MaxResults != nil {
		seen.maxResults = *params.MaxResults
	}
	f.anomaliesSeenParams = append(f.anomaliesSeenParams, seen)
	if params.NextPageToken == nil {
		f.anomaliesSeenToken = append(f.anomaliesSeenToken, "")
	} else {
		f.anomaliesSeenToken = append(f.anomaliesSeenToken, *params.NextPageToken)
	}

	if f.anomaliesErr != nil {
		return nil, f.anomaliesErr
	}
	if f.anomaliesCalls <= len(f.anomaliesPages) {
		return f.anomaliesPages[f.anomaliesCalls-1], nil
	}
	// Beyond the canned pages, keep handing out a token forever.
	return &costexplorer.GetAnomaliesOutput{
		NextPageToken: aws.String("endless"),
	}, nil
}

func anomalyPage(token string, ids ...string) *costexplorer.GetAnomaliesOutput {
	anomalies := make([]types.Anomaly, 0, len(ids))
	for _, id := range ids {
		anomalies = append(anomalies, types.Anomaly{
			AnomalyId:  aws.String(id),
			MonitorArn: aws.String("arn:aws:ce::123456789012:anomalymonitor/monitor-id"),
			AnomalyScore: &types.AnomalyScore{
				CurrentScore: 50,
				MaxScore:     50,
			},
			Impact: &types.Impact{
				MaxImpact: 100,
			},
		})
	}
	out := &costexplorer.GetAnomaliesOutput{Anomalies: anomalies}
	if token != "" {
		out.NextPageToken = aws.String(token)
	}
	return out
}

func anomaliesReq() types2.GetAnomaliesRequestType {
	return types2.GetAnomaliesRequestType{
		Time: types2.Time{Start: "2024-01-01", End: "2024-02-01"},
	}
}

func page(start string, token string, groups ...string) *costexplorer.GetCostAndUsageOutput {
	g := make([]types.Group, 0, len(groups))
	for i, name := range groups {
		g = append(g, types.Group{
			Keys: []string{name},
			Metrics: map[string]types.MetricValue{
				"UnblendedCost": {
					Amount: aws.String(strconv.Itoa(i + 1)),
					Unit:   aws.String("USD"),
				},
			},
		})
	}
	out := &costexplorer.GetCostAndUsageOutput{
		ResultsByTime: []types.ResultByTime{{
			TimePeriod: &types.DateInterval{
				Start: aws.String(start),
				End:   aws.String("2024-02-01"),
			},
			Groups: g,
		}},
	}
	if token != "" {
		out.NextPageToken = aws.String(token)
	}
	return out
}

func req() types2.CostAndUsageRequestType {
	return types2.CostAndUsageRequestType{
		Granularity: "MONTHLY",
		Metrics:     []string{"UnblendedCost"},
		Time:        types2.Time{Start: "2024-01-01", End: "2024-02-01"},
		GroupBy:     []string{"SERVICE"},
	}
}

func TestGetCostAndUsageFollowsPageTokens(t *testing.T) {
	fake := &fakeCostExplorer{pages: []*costexplorer.GetCostAndUsageOutput{
		page("2024-01-01", "t1", "EC2", "S3"),
		page("2024-01-01", "t2", "RDS"),
		page("2024-01-01", "", "Lambda"),
	}}
	srv := &Service{Client: fake}

	got, err := srv.GetCostAndUsage(context.Background(), req())
	require.NoError(t, err)

	assert.Equal(t, 3, fake.calls, "should have followed both page tokens")
	assert.Equal(t, []string{"", "t1", "t2"}, fake.seenToken,
		"each request must carry the previous page's token")

	// Every group from every page must survive into the aggregate.
	assert.Len(t, got.ResultsByTime, 3)
	var names []string
	for _, r := range got.ResultsByTime {
		for _, g := range r.Groups {
			names = append(names, g.Keys[0])
		}
	}
	assert.ElementsMatch(t, []string{"EC2", "S3", "RDS", "Lambda"}, names)

	assert.Nil(t, got.NextPageToken,
		"aggregate result must not advertise a stale page token")
}

func TestGetCostAndUsageSinglePage(t *testing.T) {
	fake := &fakeCostExplorer{pages: []*costexplorer.GetCostAndUsageOutput{
		page("2024-01-01", "", "EC2"),
	}}
	srv := &Service{Client: fake}

	got, err := srv.GetCostAndUsage(context.Background(), req())
	require.NoError(t, err)
	assert.Equal(t, 1, fake.calls)
	assert.Len(t, got.ResultsByTime, 1)
}

// An empty-string token is as terminal as a nil one.
func TestGetCostAndUsageTreatsEmptyTokenAsEnd(t *testing.T) {
	out := page("2024-01-01", "", "EC2")
	out.NextPageToken = aws.String("")
	fake := &fakeCostExplorer{pages: []*costexplorer.GetCostAndUsageOutput{out}}
	srv := &Service{Client: fake}

	_, err := srv.GetCostAndUsage(context.Background(), req())
	require.NoError(t, err)
	assert.Equal(t, 1, fake.calls, "empty token must not trigger another request")
}

func TestGetCostAndUsageStopsAtMaxPages(t *testing.T) {
	fake := &fakeCostExplorer{} // always returns a token
	srv := &Service{Client: fake}

	_, err := srv.GetCostAndUsage(context.Background(), req())
	require.Error(t, err, "an endless token stream must not loop forever")
	assert.Equal(t, maxPages, fake.calls)
	assert.Contains(t, err.Error(), "exceeded")
}

func TestGetCostAndUsagePropagatesError(t *testing.T) {
	fake := &fakeCostExplorer{err: errors.New("access denied")}
	srv := &Service{Client: fake}

	_, err := srv.GetCostAndUsage(context.Background(), req())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

// The caller's context must reach the SDK; it used to be context.TODO().
func TestGetCostAndUsageHonoursContext(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "carried")

	fake := &fakeCostExplorer{pages: []*costexplorer.GetCostAndUsageOutput{
		page("2024-01-01", "t1", "EC2"),
		page("2024-01-01", "", "S3"),
	}}
	srv := &Service{Client: fake}

	_, err := srv.GetCostAndUsage(ctx, req())
	require.NoError(t, err)
	require.Len(t, fake.seenCtx, 2)
	for i, seen := range fake.seenCtx {
		assert.Equal(t, "carried", seen.Value(ctxKey{}),
			"page %d did not receive the caller's context", i+1)
	}
}

func TestGetAnomaliesFollowsPageTokens(t *testing.T) {
	fake := &fakeCostExplorer{anomaliesPages: []*costexplorer.GetAnomaliesOutput{
		anomalyPage("t1", "a1", "a2"),
		anomalyPage("", "a3"),
	}}
	srv := &Service{Client: fake}

	got, err := srv.GetAnomalies(context.Background(), anomaliesReq())
	require.NoError(t, err)

	assert.Equal(t, 2, fake.anomaliesCalls)
	assert.Equal(t, []string{"", "t1"}, fake.anomaliesSeenToken)
	assert.Len(t, got.Anomalies, 3)
	assert.Nil(t, got.NextPageToken,
		"aggregate result must not advertise a stale page token")
}

func TestGetAnomaliesSinglePage(t *testing.T) {
	fake := &fakeCostExplorer{anomaliesPages: []*costexplorer.GetAnomaliesOutput{
		anomalyPage("", "a1"),
	}}
	srv := &Service{Client: fake}

	got, err := srv.GetAnomalies(context.Background(), anomaliesReq())
	require.NoError(t, err)
	assert.Equal(t, 1, fake.anomaliesCalls)
	assert.Len(t, got.Anomalies, 1)
}

func TestGetAnomaliesStopsAtMaxPages(t *testing.T) {
	fake := &fakeCostExplorer{} // always returns a token
	srv := &Service{Client: fake}

	_, err := srv.GetAnomalies(context.Background(), anomaliesReq())
	require.Error(t, err, "an endless token stream must not loop forever")
	assert.Equal(t, maxPages, fake.anomaliesCalls)
	assert.Contains(t, err.Error(), "exceeded")
}

func TestGetAnomaliesPropagatesError(t *testing.T) {
	fake := &fakeCostExplorer{anomaliesErr: errors.New("access denied")}
	srv := &Service{Client: fake}

	_, err := srv.GetAnomalies(context.Background(), anomaliesReq())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestGetAnomaliesAppliesOptionalFilters(t *testing.T) {
	fake := &fakeCostExplorer{anomaliesPages: []*costexplorer.GetAnomaliesOutput{
		anomalyPage("", "a1"),
	}}
	srv := &Service{Client: fake}

	req := anomaliesReq()
	req.MonitorArn = "arn:aws:ce::123456789012:anomalymonitor/monitor-id"
	req.Feedback = "YES"
	req.MaxResults = 10

	_, err := srv.GetAnomalies(context.Background(), req)
	require.NoError(t, err)

	require.Len(t, fake.anomaliesSeenParams, 1)
	sent := fake.anomaliesSeenParams[0]
	assert.Equal(t, req.MonitorArn, sent.monitorArn)
	assert.Equal(t, types.AnomalyFeedbackType(req.Feedback), sent.feedback)
	assert.Equal(t, req.MaxResults, sent.maxResults)
}

// Regression: "BundledDiscount " carried a trailing space and so never
// matched, silently leaving bundled discounts in a -l query.
func TestExcludeDiscountsRecordTypesAreNotPadded(t *testing.T) {
	expr := CostAndUsageFilterGenerator(types2.CostAndUsageRequestType{
		ExcludeDiscounts: true,
	})
	require.NotNil(t, expr)
	require.NotNil(t, expr.Not)
	require.NotNil(t, expr.Not.Dimensions)

	values := expr.Not.Dimensions.Values
	assert.Contains(t, values, "BundledDiscount")
	for _, v := range values {
		assert.Equal(t, v, strings.TrimSpace(v),
			"RECORD_TYPE value %q has surrounding whitespace and will never match", v)
	}
}
