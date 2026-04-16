package ports

import "github.com/kaiser-shaft/loyality_engine/internal/domain/models"

// порт для любого правила скидки
type DiscountRule interface {
	Name() string
	Apply(items []models.Item, currentTotal float64) float64
}

type Repository interface {
	Save(calc models.Calculation) error
	GetAll() ([]models.Calculation, error)
}
