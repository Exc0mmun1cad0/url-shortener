package httpapp

import (
	"log/slog"
	"url-shortener/internal/http/handlers/link/get"
	"url-shortener/internal/http/handlers/link/save"

	"github.com/go-chi/chi/v5"
)

func addRoutes(
	router *chi.Mux,
	log *slog.Logger,
	linkStorage LinkStorage,
) {
	router.Post("/url", save.New(log, linkStorage))
	router.Get("/url/{alias}", get.New(log, linkStorage))
}
