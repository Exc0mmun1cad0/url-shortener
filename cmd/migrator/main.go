package main

import (
	"errors"
	"flag"
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var host, port, user, password, dbName, sslMode string
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")

	flag.StringVar(&host, "host", "localhost", "postgres db server host")
	flag.StringVar(&port, "port", "5432", "postgres db server port")
	flag.StringVar(&user, "user", "", "user for postgres db server")
	flag.StringVar(&password, "password", "", "password for postgres db user")
	flag.StringVar(&dbName, "db-name", "postgres", "database name on postgser server")
	flag.StringVar(&sslMode, "ssl-mode", "require", "ssl mode for connection to postgres db")

	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is not set")
	}
	if user == "" {
		panic("postgres user is not set")
	}
	if password == "" {
		panic("password for postgres user is not set")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		postgres.FormConnStr(config.Postgres{}, host, port, user, password, dbName, sslMode, migrationsTable),
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
