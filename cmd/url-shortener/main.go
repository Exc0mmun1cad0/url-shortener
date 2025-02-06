package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	log.Info("initializing app")
	app := app.NewApp(cfg, log, storage)

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app.MustRun()
	}()
	log.Info("started app")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Stop(ctx)

	log.Info("app stopped")
}
