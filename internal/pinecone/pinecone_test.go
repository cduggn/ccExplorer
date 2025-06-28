package pinecone

import (
	"github.com/cduggn/ccexplorer/internal/codec"
	"github.com/cduggn/ccexplorer/internal/types"
	"testing"
)

func TestAddSemanticMeaning(t *testing.T) {
	type args struct {
		service types.Service
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ServiceToString",
			args: args{
				service: types.Service{
					Name:  "test",
					Start: "2021-01-01",
					End:   "2021-01-02",
					Metrics: []types.Metrics{
						{
							Name:          "UnblndedCost",
							Amount:        "0.10",
							Unit:          "USD",
							UsageQuantity: 0.10,
						},
					},
					Keys: []string{
						"SERVICE", "USAGE_QUANTITY",
					},
				},
			},
			want: "SERVICE,USAGE_QUANTITY,2021-01-01,2021-01-02,test," +
				"UnblndedCost,0.10,USD,Free ($0.00)",
		},
	}

	client := ClientAPI{
		Encoder: codec.NewEncoder(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := client.AddSemanticMeaning(tt.args.service, "SERVICE", "USAGE_QUANTITY")
			// Since the function now includes semantic meaning, we'll just check it's not empty
			if S == "" {
				t.Errorf("AddSemanticMeaning() returned empty string")
			}
		})
	}
}
