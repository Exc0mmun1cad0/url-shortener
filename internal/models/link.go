package models

import "time"

// Link is a model for url representation in storage
type Link struct {
	ID        uint       `db:"link_id" json:"id"`
	Alias     string    `db:"alias" json:"alias"`
	RawURL    string    `db:"raw_url" json:"raw_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
