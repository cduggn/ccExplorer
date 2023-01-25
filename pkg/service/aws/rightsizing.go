package aws

//type GetRightsizingRecommendationAPI interface {
//	GetRightsizingRecommendation(ctx context.Context, params *costexplorer.GetRightsizingRecommendationInput,
//		optFns ...func(*costexplorer.Options)) (*costexplorer.GetRightsizingRecommendationOutput, error)
//}

//func RightSizingRecommendationS3(ctx context.Context,
//	api GetRightsizingRecommendationAPI) (*costexplorer.
//	GetRightsizingRecommendationOutput,
//	error) {
//
//	params := &costexplorer.GetRightsizingRecommendationInput{
//		Service: aws.String("Amazon Simple Storage Service"),
//		Filter: &types.Expression{
//			Dimensions: &types.DimensionValues{
//				Key:    "SERVICE",
//				Values: []string{"Amazon Simple Storage Service"},
//			},
//		},
//	}
//
//	resp, err := api.GetRightsizingRecommendation(ctx, params)
//	if err != nil {
//		return &costexplorer.GetRightsizingRecommendationOutput{}, err
//	}
//
//	return resp, nil
//}
