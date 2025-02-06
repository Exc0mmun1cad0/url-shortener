package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

// GetURL returns the url according to its alias.
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

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

// GetLink provides information about shortened url by its alias.
func (s *Storage) GetLink(alias string) (models.Link, error) {
	const op = "storage.postgres.GetLink"

	var link models.Link
	err := s.db.Get(
		&link,
		`SELECT link_id, alias, raw_url, created_at FROM links WHERE alias = $1`,
		alias,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Link{}, storage.ErrLinkNotFound
		}
		return models.Link{}, fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}

// SaveLink adds new url shortening for further use by GetURL.
func (s *Storage) SaveLink(link models.Link) (*models.Link, error) {
	const op = "storage.postgres.CreateLink"

	var linkID uint
	err := s.db.Get(
		&linkID,
		`INSERT INTO links (alias, raw_url) VALUES ($1, $2) RETURNING link_id`,
		link.Alias, link.RawURL,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrLinkExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	link.ID = linkID
	return &link, nil
}

// DeleteLink deletes infromation about url shortening so it can be used anymore.
func (s *Storage) DeleteLink(alias string) error {
	const op = "storage.postgres.DeleteLink"

	_, err := s.db.Exec(
		`DELETE FROM links WHERE alias = $1`,
		alias,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrLinkNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
