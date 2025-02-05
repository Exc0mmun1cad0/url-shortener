package main

import (
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err.Error())
	}
	defer storage.Close()
	_ = storage
}
