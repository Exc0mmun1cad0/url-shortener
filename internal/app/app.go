package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	httpapp "url-shortener/internal/app/http"
	"url-shortener/internal/cache/redis"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
)

type App struct {
	HTTPServer *http.Server
	Storage    *postgres.Storage
	Cache      *redis.Cache
}

func NewApp(
	cfg *config.Config,
	log *slog.Logger,
	storage *postgres.Storage,
	cache *redis.Cache,
) *App {
	handler := httpapp.NewRouter(log, storage, cache)

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
		Storage:    storage,
	}
}

func (a *App) MustRun() {
	if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %s", err.Error())
	}

	if err := a.Storage.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %s", err.Error())
	}

	return nil
}
