package main

import (
	"errors"
	"flag"
	"fmt"
	"url-shortener/internal/storage/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")

	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is not set")
	}

	cfg := postgres.MustLoad()

	fmt.Println(
		postgres.FormConnStr(cfg).WithMigrationsTable(migrationsTable).String(),
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		postgres.FormConnStr(cfg).WithMigrationsTable(migrationsTable).String(),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("all migrations applied")
}
