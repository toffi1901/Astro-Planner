package entities

import "time"

type File struct {
	ID           int64     `json:"file_id"`
	UserID       *string   `json:"user_id"`
	ObjectKey    string    `json:"object_key"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}
