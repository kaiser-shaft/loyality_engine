package rules

import "github.com/kaiser-shaft/loyality_engine/internal/domain/models"

type ThresholdRule struct {
	Limit float64
	Rate  float64
}

func (r *ThresholdRule) Name() string { return "Big Receipt 5%" }
func (r *ThresholdRule) Apply(_ []models.Item, currentTotal float64) float64 {
	if currentTotal > r.Limit {
		return currentTotal * r.Rate
	}
	return 0
}
