package models

import (
	"errors"
	"strings"
	"time"
)

type Paper struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Venue     string    `json:"venue"`
	Year      int       `json:"year"`
	Status    string    `json:"status"`
	Tags      string    `json:"tags"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Paper) Validate() error {
	p.Title = strings.TrimSpace(p.Title)
	p.Venue = strings.TrimSpace(p.Venue)
	p.Status = strings.TrimSpace(strings.ToLower(p.Status))
	p.Tags = strings.TrimSpace(p.Tags)

	if p.Title == "" {
		return errors.New("title required")
	}
	if p.Year < 1900 || p.Year > time.Now().Year()+2 {
		return errors.New("year out of range")
	}
	switch p.Status {
	case "idea", "draft", "submitted", "revision", "accepted", "rejected", "published":
	default:
		return errors.New("invalid status")
	}
	return nil
}
