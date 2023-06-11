package encoder

import (
	"fmt"
	"math"
)

type Encode interface {
	CategorizeCostsWithBinning(cost float64) string
}

type Encoder struct {
}

func NewEncoder() Encode {
	return &Encoder{}
}

func (e *Encoder) CategorizeCostsWithBinning(cost float64) string {

	binRanges := []float64{0, 0.01, 1, 10, 50, 100, 500, 1000, math.MaxFloat32}
	binNames := []string{"Free", "Very Low", "Low", "Moderate", "High",
		"Very High", "Expensive"}

	binIndex := -1

	for idx, upperBound := range binRanges {
		if cost <= upperBound {
			binIndex = idx
			break
		}
	}

	var costCategory string
	if binIndex == 0 {
		costCategory = fmt.Sprintf("%s ($%.2f)", binNames[0], cost)
	} else {
		costCategory = fmt.Sprintf("%s ($%.2f)", binNames[binIndex-1], cost)
	}
	return costCategory
}
