package models

import "time"

type Image struct {
	ID          int64                  `json:"id"`
	Filename    string                 `json:"filename"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}
