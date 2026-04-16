package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kaiser-shaft/loyality_engine/internal/models"
	"github.com/kaiser-shaft/loyality_engine/pkg/render"
)

type DiscountService interface {
	CalculateAndSave(items []models.Item) (models.Calculation, error)
	GetHistory() ([]models.Calculation, error)
}

type CalculateRequest struct {
}

type Handler struct {
	service DiscountService
}

func NewHandler(service DiscountService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []models.Item `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.CalculateAndSave(req.Items)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidItem), errors.Is(err, models.ErrEmptyCart):
			render.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			render.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	render.JSON(w, result, http.StatusOK)
}

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.GetHistory()
	if err != nil {
		render.Error(w, "failed to get history", http.StatusInternalServerError)
		return
	}
	render.JSON(w, history, http.StatusOK)
}
