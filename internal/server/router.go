package server

import (
	"github.com/LealKevin/list-compare/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/", handlers.HomePageHandler)
	mux.Post("/compare", handlers.CompareHandler)

	return mux
}
