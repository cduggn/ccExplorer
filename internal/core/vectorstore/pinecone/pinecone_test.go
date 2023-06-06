package pinecone

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
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
							Name:          "test",
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
			want: "SERVICE,USAGE_QUANTITY,2021-01-01,2021-01-02,test,0.10," +
				"USD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := ServiceToString(tt.args.service)
			if S != tt.want {
				t.Errorf("ServiceToString() Got: %v, want: %v", S, tt.want)
			}
		})
	}
}
