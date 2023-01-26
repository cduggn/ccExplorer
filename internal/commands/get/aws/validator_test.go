package aws

//test ValidateGroupByMap
//func TestValidateGroupByMap(t *testing.T) {
//	type args struct {
//		groupBy map[string]string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []string
//		wantErr bool
//	}{
//		{
//			name: "Valid",
//			args: args{
//				groupBy: map[string]string{
//					"DIMENSION": "SERVICE",
//					"TAG":       "ApplicationName",
//				},
//			},
//			want:    nil,
//			wantErr: false,
//		},
//		{
//			name: "Invalid",
//			args: args{
//				groupBy: map[string]string{
//					"DIMENSION":        "SERVICE",
//					"DIMENSIONINVALID": "OPERATION",
//				},
//			},
//			want:    nil,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ValidateGroupByMap(tt.args.groupBy)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ValidateGroupByMap() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ValidateGroupByMap() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
