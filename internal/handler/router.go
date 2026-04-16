package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)
	mux.HandleFunc("GET /api/v1/history", h.History)
}
