package models

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	Title       string    `json:"title" db:"title"`
	Clicks      int       `json:"clicks" db:"clicks"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" binding:"omitempty,min=3,max=20,alphanum"`
	Title       string `json:"title,omitempty" binding:"omitempty,max=255"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type UpdateLinkRequest struct {
	OriginalURL string `json:"original_url,omitempty" binding:"omitempty,url"`
	CustomAlias string `json:"custom_alias,omitempty" binding:"omitempty,min=3,max=20,alphanum"`
	Title       string `json:"title,omitempty" binding:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type LinkResponse struct {
	ID          uuid.UUID `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	Title       string    `json:"title"`
	Clicks      int       `json:"clicks"`
	IsActive    bool      `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LinkStats struct {
	TotalLinks    int `json:"total_links"`
	TotalClicks   int `json:"total_clicks"`
	ActiveLinks   int `json:"active_links"`
	ExpiredLinks  int `json:"expired_links"`
}
