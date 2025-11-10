package models

import "time"

// ShortenRequest for POST body
type ShortenRequest struct {
	URL string `json:"url" validate:"nonzero"`
}

// ShortenResponse
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// URL entity
type URL struct {
	Shortcode string
	LongURL   string
	CreatedAt time.Time
}
