package postgres

import (
	"fmt"
	"url-shortener/internal/models"
)

// GetURL returns the url according to its alias.
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var rawURL string
	err := s.db.Select(
		&rawURL,
		`SELECT raw_url FROM links WHERE alias = $1`,
		rawURL,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return rawURL, nil
}

func (s *Storage) GetLink(alias string) (models.Link, error) {
	const op = "storage.postgres.GetLink"

	var link models.Link
	err := s.db.Get(
		&link,
		`SELECT link_id, alias, raw_url, created_at FROM links WHERE alias = $1`,
		alias,
	)
	if err != nil {
		return models.Link{}, fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}

func (s *Storage) CreateLink(link models.Link) (uint, error) {
	const op = "storage.postgres.CreateLink"

	var linkID uint
	err := s.db.Get(
		&linkID,
		`INSERT INTO links (alias, raw_url) VALUES ($1, $2) RETURNING link_id`,
		link.Alias, link.RawURL,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return linkID, nil
}

func (s *Storage) DeleteLink(alias string) (uint, error) {
	const op = "storage.postgres.DeleteLink"

	var linkID uint
	err := s.db.Get(
		&linkID,
		`DELETE FROM links WHERE alias = $1 RETURNING link_id`,
		alias,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return linkID, nil
}
