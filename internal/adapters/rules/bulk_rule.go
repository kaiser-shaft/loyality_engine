package rules

import "github.com/kaiser-shaft/loyality_engine/internal/domain/models"

type BulkRule struct{ N int }

func (r *BulkRule) Name() string { return "Bulk 3+1" }
func (r *BulkRule) Apply(items []models.Item, _ float64) float64 {
	var sum float64
	for _, it := range items {
		if it.Quantity >= r.N {
			sum += float64(it.Quantity/r.N) * it.Price
		}
	}
	return sum
}
