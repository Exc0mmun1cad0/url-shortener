package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/model"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

// GetURL returns the url according to its alias.
func (s *Storage) GetURLByAlias(alias string) (string, error) {
	const op = "storage.postgres.GetAlias"

	var rawURL string
	err := s.db.Get(
		&rawURL,
		`SELECT raw_url FROM links WHERE alias = $1`,
		alias,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return rawURL, nil
}

// GetURL provides information about shortened url by its alias.
func (s *Storage) GetURL(alias string) (models.URL, error) {
	const op = "storage.postgres.GetURL"

	var url models.URL
	err := s.db.Get(
		&url,
		`SELECT link_id, alias, raw_url, created_at FROM links WHERE alias = $1`,
		alias,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.URL{}, storage.ErrURLNotFound
		}
		return models.URL{}, fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

// SaveURL adds new url shortening for further use by GetURL.
func (s *Storage) SaveURL(url models.URL) (*models.URL, error) {
	const op = "storage.postgres.CreateURL"

	var newURL models.URL
	err := s.db.Get(
		&newURL,
		`INSERT INTO links (alias, raw_url) VALUES ($1, $2) RETURNING *`,
		url.Alias, url.RawURL,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &newURL, nil
}

// DeleteURL deletes infromation about url shortening so it can be used anymore.
func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	res, err := s.db.Exec(
		`DELETE FROM links WHERE alias = $1`,
		alias,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Check whether deletion was successful
	if num, _ := res.RowsAffected(); num == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
