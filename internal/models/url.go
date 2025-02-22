package models

import "time"

// URL is a model for url representation in storage
type URL struct {
	ID        uint      `db:"link_id" json:"id"`
	Alias     string    `db:"alias" json:"alias"`
	RawURL    string    `db:"raw_url" json:"raw_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
