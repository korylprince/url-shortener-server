package db

import "time"

//URL represents a shortened URL
type URL struct {
	URL     string     `json:"url"`
	Views   uint64     `json:"views"`
	Expires *time.Time `json:"expires"`
}
