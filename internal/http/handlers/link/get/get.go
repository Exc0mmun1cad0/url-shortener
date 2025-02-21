package get

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Response resp.Response
	Link     models.Link `json:"link"`
}

type LinkGetter interface {
	GetLink(alias string) (models.Link, error)
}

// New creates handler for requests connected with link information.
func New(log *slog.Logger, linkGetter LinkGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.link.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("reques_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		link, err := linkGetter.GetLink(alias)
		if errors.Is(err, storage.ErrLinkNotFound) {
			log.Info("alias information not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("alias information not found"))

			return
		}
		if err != nil {
			log.Error("failed to get alias information")

			render.JSON(w, r, resp.Error("failed to get alias information"))

			return
		}

		log.Info("alias information found", slog.Any("link", link))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Link:     link,
		})
	}
}
