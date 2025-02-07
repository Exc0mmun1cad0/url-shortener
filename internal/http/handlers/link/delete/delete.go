package delete

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/cache"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type LinkDeleter interface {
	DeleteLink(alias string) error
}

type CacheDeleter interface {
	Delete(ctx context.Context, key string) error
}

func New(log *slog.Logger, linkDeleter LinkDeleter, cacheDeleter CacheDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.link.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		err := linkDeleter.DeleteLink(alias)
		if errors.Is(err, storage.ErrLinkNotFound) {
			log.Info("link not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("link not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete link by its alias")

			render.JSON(w, r, resp.Error("failed to delete link by its alias"))

			return
		}

		log.Info("link deleted", slog.String("alias", alias))

		err = cacheDeleter.Delete(r.Context(), alias)
		if errors.Is(err, cache.ErrNotFound) {
			slog.Info("url not found in cache", slog.String("alias", alias))
		}

		log.Info("link deleted fron cache", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())
	}
}
