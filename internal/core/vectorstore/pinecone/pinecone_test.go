package pinecone

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/encoder"
	"testing"
)

func TestServiceToString(t *testing.T) {
	type args struct {
		service model.Service
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ServiceToString",
			args: args{
				service: model.Service{
					Name:  "test",
					Start: "2021-01-01",
					End:   "2021-01-02",
					Metrics: []model.Metrics{
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
		Encoder: encoder.NewEncoder(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := client.serviceToString(tt.args.service)
			if S != tt.want {
				t.Errorf("ServiceToString() Got: %v, want: %v", S, tt.want)
			}
		})
	}
}
