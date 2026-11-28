package model

import (
	"time"

	"uuid"
)

type Track struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Album       string    `json:"album"`
	ObjectKey   string    `json:"object_key"`
	Duration    int       `json:"duration"`
	MIMEType    string    `json:"mime_type"`
	FileSize    int64     `json:"file_size"`
	CoverArtURL string    `json:"cover_art_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
