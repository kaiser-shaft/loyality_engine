package service

import (
	"math"

	"github.com/kaiser-shaft/loyality_engine/internal/models"
)

const (
	foodDiscountRate    = 0.1
	receiptDiscountRate = 0.05
	everyNthItemFree    = 4
	thresholdAmount     = 10000.0
)

type CalculationRepository interface {
	Save(calc models.Calculation) error
	GetAll() ([]models.Calculation, error)
}

type discountService struct {
	repo CalculationRepository
}

func NewDiscountService(repo CalculationRepository) *discountService {
	return &discountService{repo: repo}
}

func (s *discountService) CalculateAndSave(items []models.Item) (models.Calculation, error) {
	if len(items) == 0 {
		return models.Calculation{}, models.ErrEmptyCart
	}
	for _, it := range items {
		if it.Price <= 0 || it.Quantity <= 0 || it.Name == "" || it.ID == "" {
			return models.Calculation{}, models.ErrInvalidItem
		}
	}
	calc := s.calculate(items)
	if err := s.repo.Save(calc); err != nil {
		return models.Calculation{}, models.ErrInternalStorage
	}
	return calc, nil
}

func (s *discountService) calculate(items []models.Item) models.Calculation {
	var originalTotal float64
	var discounts []models.AppliedDiscount

	for _, item := range items {
		originalTotal += item.Price * float64(item.Quantity)
	}

	currentTotal := originalTotal

	// Правило 1: Скидка 10% на категорию food
	var foodDiscountSum float64
	for _, item := range items {
		if item.Category == "food" {
			foodDiscountSum += (item.Price * float64(item.Quantity)) * foodDiscountRate
		}
	}
	if foodDiscountSum > 0 {
		foodDiscountSum = round(foodDiscountSum)
		discounts = append(discounts, models.AppliedDiscount{Name: "Food Category 10%", Amount: foodDiscountSum})
		currentTotal -= foodDiscountSum
	}

	// Правило 2: Опт 3+1 (каждый 4-й бесплатно)
	var bulkDiscountSum float64
	for _, item := range items {
		if item.Quantity >= everyNthItemFree {
			freeItems := item.Quantity / everyNthItemFree
			bulkDiscountSum += float64(freeItems) * item.Price
		}
	}
	if bulkDiscountSum > 0 {
		bulkDiscountSum = round(bulkDiscountSum)
		discounts = append(discounts, models.AppliedDiscount{Name: "Bulk 3+1", Amount: bulkDiscountSum})
		currentTotal -= bulkDiscountSum
	}

	// Правило 3: Скидка за объем чека (от текущего остатка)
	if currentTotal > thresholdAmount {
		receiptDiscount := currentTotal * receiptDiscountRate
		receiptDiscount = round(receiptDiscount)
		discounts = append(discounts, models.AppliedDiscount{Name: "Big Receipt 5%", Amount: receiptDiscount})
		currentTotal -= receiptDiscount
	}

	return models.Calculation{
		OriginalTotal: round(originalTotal),
		Discounts:     discounts,
		FinalTotal:    round(currentTotal),
	}
}

func (s *discountService) SaveCalculation(calc models.Calculation) error {
	return s.repo.Save(calc)
}

func (s *discountService) GetHistory() ([]models.Calculation, error) {
	return s.repo.GetAll()
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}
