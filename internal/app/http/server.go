package httpapp

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/handlers/link/save"
	mw "url-shortener/internal/http/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates router for our app. It will be used in
// creating http.Server as a handler.
func NewRouter(
	log *slog.Logger,
	linkSaver save.LinkSaver,
) http.Handler {
	router := chi.NewRouter()

	// Add middlewares for router
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		mw.NewLogger(log),
		middleware.Recoverer,
		middleware.URLFormat,
	)

	addRoutes(router, log, linkSaver)

	return router
}
