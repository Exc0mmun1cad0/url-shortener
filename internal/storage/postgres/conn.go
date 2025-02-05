package postgres

import (
	"fmt"
	"net/url"
	"url-shortener/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates new postgres connection.
func New(cfg config.Postgres) (*Storage, error) {
	const op = "storage.postgres.New"

	// construct string for connection
	q := url.Values{}
	q.Set("sslmode", cfg.SSLMode)
	dataSrc := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:     cfg.DBName,
		RawQuery: q.Encode(),
	}

	// establish a connection
	conn, err := sqlx.Open("postgres", dataSrc.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// check connection by doing ping
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

// Close closes postgres connection.
func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
