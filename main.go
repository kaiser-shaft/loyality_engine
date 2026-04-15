package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
)

// Константы для бизнес-логики
const (
	foodDiscountRate    = 0.1
	receiptDiscountRate = 0.05
	everyNthItemFree    = 4
	thresholdAmount     = 10000.0
)

// Структуры данных
type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Quantity int     `json:"quantity"`
}

type CalculateRequest struct {
	Items []Item `json:"items"`
}

type AppliedDiscount struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type CalculateResponse struct {
	OriginalTotal float64           `json:"original_total"`
	Discounts     []AppliedDiscount `json:"discounts"`
	FinalTotal    float64           `json:"final_total"`
}

func main() {
	// Настройка логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	// Используем новый синтаксис роутинга Go 1.22+
	mux.HandleFunc("POST /api/v1/calculate", calculateHandler)

	addr := ":8080"
	slog.Info("server started", "addr", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// 1. Валидация
	if len(req.Items) == 0 {
		renderError(w, "items list is empty", http.StatusBadRequest)
		return
	}

	for _, item := range req.Items {
		if item.Price <= 0 || item.Quantity <= 0 || item.ID == "" || item.Name == "" {
			renderError(w, fmt.Sprintf("invalid item: %s", item.Name), http.StatusBadRequest)
			return
		}
	}

	var originalTotal float64
	var discounts []AppliedDiscount

	// Считаем общую сумму без скидок
	for _, item := range req.Items {
		originalTotal += item.Price * float64(item.Quantity)
	}

	currentTotal := originalTotal

	// 2. Логика скидок (Применяем ПОСЛЕДОВАТЕЛЬНО)

	// Правило 1: Скидка 10% на категорию food
	var foodDiscountSum float64
	for _, item := range req.Items {
		if item.Category == "food" {
			foodDiscountSum += (item.Price * float64(item.Quantity)) * foodDiscountRate
		}
	}
	if foodDiscountSum > 0 {
		foodDiscountSum = round(foodDiscountSum)
		discounts = append(discounts, AppliedDiscount{Name: "Food Category 10%", Amount: foodDiscountSum})
		currentTotal -= foodDiscountSum
	}

	// Правило 2: Опт 3+1 (каждый 4-й бесплатно)
	var bulkDiscountSum float64
	for _, item := range req.Items {
		if item.Quantity >= everyNthItemFree {
			freeItems := item.Quantity / everyNthItemFree
			bulkDiscountSum += float64(freeItems) * item.Price
		}
	}
	if bulkDiscountSum > 0 {
		bulkDiscountSum = round(bulkDiscountSum)
		discounts = append(discounts, AppliedDiscount{Name: "Bulk 3+1", Amount: bulkDiscountSum})
		currentTotal -= bulkDiscountSum
	}

	// Правило 3: Скидка за объем чека (от текущего остатка)
	if currentTotal > thresholdAmount {
		receiptDiscount := currentTotal * receiptDiscountRate
		receiptDiscount = round(receiptDiscount)
		discounts = append(discounts, AppliedDiscount{Name: "Big Receipt 5%", Amount: receiptDiscount})
		currentTotal -= receiptDiscount
	}

	// 3. Формирование ответа
	resp := CalculateResponse{
		OriginalTotal: round(originalTotal),
		Discounts:     discounts,
		FinalTotal:    round(currentTotal),
	}

	renderJSON(w, resp, http.StatusOK)
}

// Вспомогательные функции
func round(val float64) float64 {
	return math.Round(val*100) / 100
}

func renderJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func renderError(w http.ResponseWriter, msg string, status int) {
	renderJSON(w, map[string]string{"error": msg}, status)
}
