package httpapp

import (
	"log/slog"
	"net/http"
	del "url-shortener/internal/http/handlers/url/delete"
	"url-shortener/internal/http/handlers/url/get"
	"url-shortener/internal/http/handlers/url/save"
	"url-shortener/internal/http/handlers/redirect"
	mw "url-shortener/internal/http/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type URLStorage interface {
	save.URLSaver
	get.URLGetter
	del.URLDeleter

	redirect.URLByAliasGetter
}

type URLCache interface {
	redirect.CacheInsertGetter
	del.CacheDeleter
}

// NewRouter creates router for our app. It will be used in
// creating http.Server as a handler.
func NewRouter(
	log *slog.Logger,
	urlStorage URLStorage,
	urlCache URLCache,
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

	addRoutes(router, log, urlStorage, urlCache)

	return router
}
