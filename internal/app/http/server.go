package httpapp

import (
	"log/slog"
	"net/http"
	del "url-shortener/internal/http/handlers/link/delete"
	"url-shortener/internal/http/handlers/link/get"
	"url-shortener/internal/http/handlers/link/save"
	"url-shortener/internal/http/handlers/redirect"
	mw "url-shortener/internal/http/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type LinkStorage interface {
	save.LinkSaver
	get.LinkGetter
	del.LinkDeleter

	redirect.URLGetter
}

type LinkCache interface {
	redirect.CacheInsertGetter
	del.CacheDeleter
}

// NewRouter creates router for our app. It will be used in
// creating http.Server as a handler.
func NewRouter(
	log *slog.Logger,
	linkStorage LinkStorage,
	linkCache LinkCache,
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

	addRoutes(router, log, linkStorage, linkCache)

	return router
}
