package logic

import (
	"math"

	"github.com/kaiser-shaft/loyality_engine/internal/domain/models"
	"github.com/kaiser-shaft/loyality_engine/internal/domain/ports"
)

type LoyaltyEngine struct {
	repo  ports.Repository
	rules []ports.DiscountRule
}

func NewLoyaltyEngine(repo ports.Repository, rules []ports.DiscountRule) *LoyaltyEngine {
	return &LoyaltyEngine{repo: repo, rules: rules}
}

func (e *LoyaltyEngine) CalculateAndSave(items []models.Item) (models.Calculation, error) {
	if len(items) == 0 {
		return models.Calculation{}, models.ErrEmptyCart
	}
	var originalTotal float64
	for _, it := range items {
		if it.Price <= 0 || it.Quantity <= 0 || it.Name == "" || it.ID == "" {
			return models.Calculation{}, models.ErrInvalidItem
		}
	}
	currentTotal := originalTotal
	var discounts []models.AppliedDiscount
	for _, rule := range e.rules {
		amount := rule.Apply(items, currentTotal)
		if amount > 0 {
			amount = round(amount)
			currentTotal -= amount
			discounts = append(discounts, models.AppliedDiscount{
				Name:   rule.Name(),
				Amount: amount,
			})
		}
	}
	calc := models.Calculation{
		OriginalTotal: originalTotal,
		Discounts:     discounts,
		FinalTotal:    currentTotal,
	}
	if err := e.repo.Save(calc); err != nil {
		return models.Calculation{}, models.ErrInternalStorage
	}
	return calc, nil
}

func (e *LoyaltyEngine) GetHistory() ([]models.Calculation, error) {
	return e.repo.GetAll()
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}
