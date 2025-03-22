package get

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/model"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Response resp.Response
	URL     models.URL `json:"url"`
}

type URLGetter interface {
	GetURL(alias string) (models.URL, error)
}

// New creates handler for requests connected with url information.
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("reques_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("alias information not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("alias information not found"))

			return
		}
		if err != nil {
			log.Error("failed to get alias information")

			render.JSON(w, r, resp.Error("failed to get alias information"))

			return
		}

		log.Info("alias information found", slog.Any("url", url))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			URL:     url,
		})
	}
}
