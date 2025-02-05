package httpapp

import (
	"log/slog"
	"url-shortener/internal/config"

	"github.com/go-chi/chi/v5"
)

func addRoutes(
	chi *chi.Mux,
	log *slog.Logger,
	config *config.Config,

) {
	// TODO: register handlers for chi router
}
