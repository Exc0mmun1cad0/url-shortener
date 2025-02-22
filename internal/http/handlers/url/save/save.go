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
	RawURL   string `json:"raw_url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Response resp.Response
	URL     models.URL `json:"url"`
}

type URLSaver interface {
	SaveURL(url models.URL) (*models.URL, error)
}

// New creates handler for creating new urls.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
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

			return
		}

		alias := req.Alias
		if alias == "" {
			alias, err = als.Generate(req.RawURL, aliasLength)
			if err != nil {
				log.Error("failed to generate alias", sl.Err(err))

				render.JSON(w, r, resp.Error("failed to generate alias"))
			}
		}

		url, err := urlSaver.SaveURL(models.URL{Alias: alias, RawURL: req.RawURL})
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.RawURL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Any("url", *url))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			URL:     *url,
		})
	}
}
