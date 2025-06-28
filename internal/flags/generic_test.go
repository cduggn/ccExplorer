package flags

import (
	"testing"
)

func TestDimensionValidator_Validate(t *testing.T) {
	validator := DimensionValidator{}
	
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid dimension",
			input:   "DIMENSION=SERVICE",
			wantErr: false,
		},
		{
			name:    "valid tag",
			input:   "TAG=Environment",
			wantErr: false,
		},
		{
			name:    "multiple valid entries",
			input:   "DIMENSION=SERVICE,TAG=Environment",
			wantErr: false,
		},
		{
			name:    "invalid dimension",
			input:   "DIMENSION=INVALID",
			wantErr: true,
		},
		{
			name:    "invalid type",
			input:   "INVALID=SERVICE",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DimensionValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify that result has some data when successful
				if len(result.Dimensions) == 0 && len(result.Tags) == 0 {
					t.Errorf("Expected some dimensions or tags to be parsed")
				}
			}
		})
	}
}

func TestFilterValidator_Validate(t *testing.T) {
	validator := FilterValidator{}
	
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid service filter",
			input:   "SERVICE=EC2-Instance",
			wantErr: false,
		},
		{
			name:    "valid tag filter",
			input:   "TAG=Production",
			wantErr: false,
		},
		{
			name:    "multiple filters",
			input:   "SERVICE=EC2-Instance,REGION=us-east-1",
			wantErr: false,
		},
		{
			name:    "invalid dimension",
			input:   "INVALID=value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify that result has some data when successful
				if len(result.Dimensions) == 0 && len(result.Tags) == 0 {
					t.Errorf("Expected some dimensions or tags to be parsed")
				}
			}
		})
	}
}

func TestGenericFlag_SetAndValue(t *testing.T) {
	flag := NewGroupByFlag()
	
	err := flag.Set("DIMENSION=SERVICE,TAG=Environment")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !flag.IsSet() {
		t.Error("Expected flag to be marked as set")
	}
	
	value := flag.Value()
	if len(value.Dimensions) != 1 || value.Dimensions[0] != "SERVICE" {
		t.Errorf("Expected dimension SERVICE, got %v", value.Dimensions)
	}
	
	if len(value.Tags) != 1 || value.Tags[0] != "Environment" {
		t.Errorf("Expected tag Environment, got %v", value.Tags)
	}
}