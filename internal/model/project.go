package model

import "time"

type Project struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	EncryptedName    []byte    `json:"encrypted_name"`
	EncryptedContent []byte    `json:"encrypted_content"`
	SortOrder        int       `json:"sort_order"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	EncryptedName    []byte `json:"encrypted_name"`
	EncryptedContent []byte `json:"encrypted_content"`
	SortOrder        int    `json:"sort_order"`
}

type UpdateProjectRequest struct {
	EncryptedName    []byte `json:"encrypted_name"`
	EncryptedContent []byte `json:"encrypted_content"`
	SortOrder        int    `json:"sort_order"`
}

type ProjectMeta struct {
	ID        int64     `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}
