package aws

//type mockGetCostAndUsageAPI func(ctx context.Context,
//	optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
//
//func (m mockGetCostAndUsageAPI) GetCostAndUsage(ctx context.Context,
//	optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {
//	return m(ctx, optFns...)
//}
//
//var mockOutput = &costexplorer.GetCostAndUsageOutput{
//	ResultsByTime: []types.ResultByTime{
//		{
//			Estimated: true,
//			TimePeriod: &types.DateInterval{
//				End:   aws.String("2020-01-01"),
//				Start: aws.String("2020-12-31"),
//			},
//			Total: map[string]types.MetricValue{
//				"UnblendedCost": {
//					Amount: aws.String("100"),
//					Unit:   aws.String("USD"),
//				},
//				"BlendedCost": {
//					Amount: aws.String("100"),
//					Unit:   aws.String("USD"),
//				},
//			},
//		},
//	},
//}

//func TestAPIClient_GetCostAndUsage(t *testing.T) {
//
//	cases := []struct {
//		client func(t *testing.T) model.GetCostAndUsageAPI
//		bucket string
//		key    string
//		expect *costexplorer.GetCostAndUsageOutput
//	}{
//		{
//			client: func(t *testing.T) model.GetCostAndUsageAPI {
//				return mockGetCostAndUsageAPI(func(ctx context.Context,
//					optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {
//					t.Helper()
//
//					return mockOutput, nil
//				})
//			},
//			expect: mockOutput,
//		},
//	}
//
//
//	for i, tt := range cases {
//		t.Run(strconv.Itoa(i), func(t *testing.T) {
//			ctx := context.TODO()
//
//			results, err := awsClient.GetCostAndUsage(ctx,
//				model.CostAndUsageRequestType{
//					Granularity: "MONTHLY",
//					Time: model.Time{
//						Start: "2020-01-01",
//						End:   "2020-12-31",
//					},
//					GroupBy: []string{"SERVICE"},
//				})
//			if err != nil {
//				t.Fatalf("expect no error, got %v", err)
//			}
//			//if e, a := tt.expect, results; results.Services != 0 {
//			//	t.Errorf("expect %v, got %v", e, a)
//			//}
//			assert.Equal(t, len(tt.expect.ResultsByTime), len(results.ResultsByTime))
//		})
//	}
//
//}

//type mockGetCostForecastAPI func(ctx context.Context,
//	params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error)
//
//func (m mockGetCostForecastAPI) GetCostForecast(ctx context.Context,
//	params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error) {
//	return m(ctx, params, optFns...)
//}
//
//var mockGetForecastOutput = &costexplorer.GetCostForecastOutput{
//	Total: &types.MetricValue{
//		Amount: aws.String("100"),
//		Unit:   aws.String("USD"),
//	},
//}

//func TestGetCostForecast(t *testing.T) {
//	cases := []struct {
//		client func(t *testing.T) model.GetCostForecastAPI
//		bucket string
//		key    string
//		expect *costexplorer.GetCostForecastOutput
//	}{
//		{
//			client: func(t *testing.T) model.GetCostForecastAPI {
//				return mockGetCostForecastAPI(func(ctx context.Context,
//					params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error) {
//					t.Helper()
//
//					return mockGetForecastOutput, nil
//				})
//			},
//			expect: mockGetForecastOutput,
//		},
//	}
//
//	for i, tt := range cases {
//		t.Run(strconv.Itoa(i), func(t *testing.T) {
//			client := tt.client(t)
//			output, err := client.GetCostForecast(context.Background(), &costexplorer.GetCostForecastInput{})
//			if err != nil {
//				t.Fatal(err)
//			}
//			if output.Total.Amount != tt.expect.Total.Amount {
//				t.Errorf("expected %v, got %v", tt.expect.Total.Amount, output.Total.Amount)
//			}
//		})
//	}
//
//}
