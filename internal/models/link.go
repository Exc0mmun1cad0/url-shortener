package models

import "time"

// Link is a model for url representation in storage
type Link struct {
	ID        uint       `db:"link_id"`
	Alias     string    `db:"alias"`
	RawURL    string    `db:"raw_url"`
	CreatedAt time.Time `db:"created_at"`
}
