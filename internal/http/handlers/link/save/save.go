package save

import (
	"errors"
	"log/slog"
	"net/http"
	als "url-shortener/internal/lib/alias"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// TODO: move to config
const (
	aliasLength = 7
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Response resp.Response
	Link     models.Link `json:"link"`
}

type linkSaver interface {
	SaveLink(link models.Link) (*models.Link, error)
}

func New(log *slog.Logger, linkSaver linkSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.link.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		defer r.Body.Close()

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
		}

		alias := req.Alias
		if alias == "" {
			alias, err = als.Generate(req.URL, aliasLength)
			if err != nil {
				log.Error("failed to generate alias", sl.Err(err))

				render.JSON(w, r, resp.Error("failed to generate alias"))
			}
		}

		link, err := linkSaver.SaveLink(models.Link{Alias: alias, RawURL: req.URL})
		if errors.Is(err, storage.ErrLinkExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		log.Info("link added", slog.Any("link", *link))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Link:     *link,
		})
	}
}
