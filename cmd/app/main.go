package main

import (
	"log/slog"
	"net/http"

	"github.com/kaiser-shaft/loyality_engine/internal/adapters/repository"
	"github.com/kaiser-shaft/loyality_engine/internal/adapters/rules"
	"github.com/kaiser-shaft/loyality_engine/internal/domain/logic"
	"github.com/kaiser-shaft/loyality_engine/internal/domain/ports"
	"github.com/kaiser-shaft/loyality_engine/internal/handler"
)

func main() {
	// 1. Инфраструктура
	repo := repository.NewInMemoryRepo()

	// 2. Выбираем активные правила (порядок важен!)
	activeRules := []ports.DiscountRule{
		&rules.FoodRule{Rate: 0.1},
		&rules.BulkRule{N: 4},
		&rules.ThresholdRule{Limit: 10000, Rate: 0.05},
	}

	// 3. Собираем ядро
	engine := logic.NewLoyaltyEngine(repo, activeRules)

	// 4. Подключаем транспортный адаптер
	h := handler.NewHandler(engine)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, h)

	slog.Info("server started on :8080")
	http.ListenAndServe(":8080", mux)
}
