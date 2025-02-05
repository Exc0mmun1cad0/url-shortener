package httpapp

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/config"

	"github.com/go-chi/chi/v5"
)

// NewRouter creates router for our app. It will be used in
// creating http.Server in handler field.
func NewRouter(
	log *slog.Logger,
	config *config.Config,
) http.Handler {
	mux := chi.NewMux()
	addRoutes(
		mux,
		log,
		config,
	)

	// TODO: add middleware registration. How should it be:
	// var handler http.Handler = mux
	// handler = someMiddleware(handler)
	// handler = someMiddleware2(handler)
	// handler = someMiddleware3(handler)
	// return handler

	return mux
}
