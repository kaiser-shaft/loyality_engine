package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/kaiser-shaft/loyality_engine/internal/handler"
	"github.com/kaiser-shaft/loyality_engine/internal/repository"
	"github.com/kaiser-shaft/loyality_engine/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo := repository.NewInMemoryRepo()
	svc := service.NewDiscountService(repo)
	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, h)

	addr := ":8080"
	slog.Info("server started", "addr", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
