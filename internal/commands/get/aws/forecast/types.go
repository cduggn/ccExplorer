package forecast

import (
	aws2 "github.com/cduggn/ccexplorer/pkg/domain/model"
)

type CommandLineInput struct {
	FilterByValues          aws2.Filter
	Granularity             string
	PredictionIntervalLevel int32
	Start                   string
	End                     string
}
