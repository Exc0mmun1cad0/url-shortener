package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/cache"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

type CacheInsertGetter interface {
	Insert(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter, cacheIO CacheInsertGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		// toAdd indicates whether i should add it (alias, url) in cache or not
		toAdd := false
		// first of all, we search it in cache
		url, err := cacheIO.Get(context.TODO(), alias)
		if err != nil {
			if errors.Is(err, cache.ErrNotFound) {
				log.Info("url not found in cache")
				toAdd = true
			} else {
				log.Error("failed to find url in cache", sl.Err(err))

			}
		}

		// if url was found, we redirect to it
		if url != "" {
			log.Info("found url in cache", slog.String("url", url))

			redirect(w, r, url)
			return
		}

		url, err = urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			log.Error("failed to find url", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("found url", slog.String("url", url))

		if toAdd {
			err = cacheIO.Insert(context.TODO(), alias, url)
			if err != nil {
				slog.Error("failed to add url to cache", sl.Err(err))
			}

			log.Info("added url to cache", slog.String("url", url))
		}

		redirect(w, r, url)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)

}
