package httpapp

import (
	"log/slog"
	"url-shortener/internal/http/handlers/link/save"

	"github.com/go-chi/chi/v5"
)

func addRoutes(
	router *chi.Mux,
	log *slog.Logger,
	linkSaver save.LinkSaver,
) {
	router.Post("/url", save.New(log, linkSaver))
}
