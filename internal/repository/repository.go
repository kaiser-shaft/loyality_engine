package repository

import (
	"sync"

	"github.com/kaiser-shaft/loyality_engine/internal/models"
)

type inMemoryRepo struct {
	sync.RWMutex
	calculations []models.Calculation
}

func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{}
}

func (r *inMemoryRepo) Save(calc models.Calculation) error {
	r.Lock()
	defer r.Unlock()

	r.calculations = append(r.calculations, calc)

	return nil
}

func (r *inMemoryRepo) GetAll() ([]models.Calculation, error) {
	r.RLock()
	defer r.RUnlock()

	return r.calculations, nil
}
