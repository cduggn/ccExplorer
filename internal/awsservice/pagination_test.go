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
