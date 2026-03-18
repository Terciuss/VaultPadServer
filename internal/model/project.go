package model

import "time"

type Project struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	EncryptedContent []byte    `json:"encrypted_content"`
	KeyCheck         []byte    `json:"key_check"`
	SortOrder        int       `json:"sort_order"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name             string `json:"name"`
	EncryptedContent []byte `json:"encrypted_content"`
	KeyCheck         []byte `json:"key_check"`
	SortOrder        int    `json:"sort_order"`
}

type UpdateProjectRequest struct {
	Name             string `json:"name"`
	EncryptedContent []byte `json:"encrypted_content"`
	KeyCheck         []byte `json:"key_check"`
	SortOrder        int    `json:"sort_order"`
}

type ProjectMeta struct {
	ID        int64     `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}
