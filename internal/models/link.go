package link

import "time"

// Link is a model for url representation in storage
type Link struct {
	Alias     string
	RawURL    string
	CreatedAt time.Time
}
