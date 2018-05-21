package db

import "time"

//URL represents a shortened URL
type URL struct {
	ID           string     `json:"id"`
	User         string     `json:"user"`
	URL          string     `json:"url"`
	Views        uint64     `json:"views"`
	Expires      *time.Time `json:"expires"`
	LastModified *time.Time `json:"last_modified"`
}
