package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates new postgres connection.
func New(cfg Config) (*Storage, error) {
	const op = "storage.postgres.New"

	// establish a connection
	conn, err := sqlx.Open("postgres", FormConnStr(cfg).String())
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
