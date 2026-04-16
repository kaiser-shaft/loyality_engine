package rules

import "github.com/kaiser-shaft/loyality_engine/internal/domain/models"

type FoodRule struct{ Rate float64 }

func (r *FoodRule) Name() string { return "Food Category 10%" }
func (r *FoodRule) Apply(items []models.Item, _ float64) float64 {
	var sum float64
	for _, it := range items {
		if it.Category == "food" {
			sum += it.Price * float64(it.Quantity) * r.Rate
		}
	}
	return sum
}
