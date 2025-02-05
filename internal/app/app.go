package app

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	httpapp "url-shortener/internal/app/http"
	"url-shortener/internal/config"
)

type App struct {
	HTTPServer *http.Server
	cfg        *config.Config
}

func NewApp(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	handler := httpapp.NewRouter(log, cfg) //TODO: maybe i should pass db connection here?

	address := net.JoinHostPort(cfg.HTTPServer.Host, strconv.Itoa(cfg.HTTPServer.Port))
	srv := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadTimeout:       cfg.HTTPServer.Timeout,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		HTTPServer: srv,
		cfg: cfg,
	}
}

func (a *App) Run() {

}
