package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		url, err := urlGetter.GetURL(alias)
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

		http.Redirect(w, r, url, http.StatusFound)
	}
}
