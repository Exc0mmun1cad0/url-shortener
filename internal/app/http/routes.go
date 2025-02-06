package httpapp

import (
	"log/slog"
	del "url-shortener/internal/http/handlers/link/delete"
	"url-shortener/internal/http/handlers/link/get"
	"url-shortener/internal/http/handlers/link/save"
	"url-shortener/internal/http/handlers/redirect"

	"github.com/go-chi/chi/v5"
)

func addRoutes(
	router *chi.Mux,
	log *slog.Logger,
	linkStorage LinkStorage,
) {
	router.Post("/url", save.New(log, linkStorage))
	router.Get("/url/{alias}", get.New(log, linkStorage))
	router.Delete("/url/{alias}", del.New(log, linkStorage))

	router.Get("/{alias}", redirect.New(log, linkStorage))
}
