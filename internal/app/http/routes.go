package httpapp

import (
	"log/slog"
	del "url-shortener/internal/http/handlers/url/delete"
	"url-shortener/internal/http/handlers/url/get"
	"url-shortener/internal/http/handlers/url/save"
	"url-shortener/internal/http/handlers/redirect"

	"github.com/go-chi/chi/v5"
)

func addRoutes(
	router *chi.Mux,
	log *slog.Logger,
	urlStorage URLStorage,
	urlCache URLCache,
) {
	router.Post("/url", save.New(log, urlStorage))
	router.Get("/url/{alias}", get.New(log, urlStorage))
	router.Delete("/url/{alias}", del.New(log, urlStorage, urlCache))

	router.Get("/{alias}", redirect.New(log, urlStorage, urlCache))
}
