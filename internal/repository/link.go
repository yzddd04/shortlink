package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"link-shortener/internal/database"
	"link-shortener/internal/models"
)

type LinkRepository struct {
	db *database.Database
}

func NewLinkRepository(db *database.Database) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(link *models.Link) error {
	query := `
		INSERT INTO links (id, user_id, original_url, short_code, title, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	
	return r.db.DB.QueryRow(
		query,
		link.ID,
		link.UserID,
		link.OriginalURL,
		link.ShortCode,
		link.Title,
		link.ExpiresAt,
	).Scan(&link.CreatedAt, &link.UpdatedAt)
}

func (r *LinkRepository) GetByID(id uuid.UUID) (*models.Link, error) {
	link := &models.Link{}
	query := `
		SELECT id, user_id, original_url, short_code, title, clicks, is_active, expires_at, created_at, updated_at
		FROM links WHERE id = $1
	`
	
	err := r.db.DB.QueryRow(query, id).Scan(
		&link.ID,
		&link.UserID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.Title,
		&link.Clicks,
		&link.IsActive,
		&link.ExpiresAt,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("link not found")
		}
		return nil, err
	}
	
	return link, nil
}

func (r *LinkRepository) GetByShortCode(shortCode string) (*models.Link, error) {
	link := &models.Link{}
	query := `
		SELECT id, user_id, original_url, short_code, title, clicks, is_active, expires_at, created_at, updated_at
		FROM links WHERE short_code = $1 AND is_active = true
	`
	
	err := r.db.DB.QueryRow(query, shortCode).Scan(
		&link.ID,
		&link.UserID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.Title,
		&link.Clicks,
		&link.IsActive,
		&link.ExpiresAt,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("link not found")
		}
		return nil, err
	}
	
	// Check if link is expired
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, fmt.Errorf("link has expired")
	}
	
	return link, nil
}

func (r *LinkRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Link, error) {
	query := `
		SELECT id, user_id, original_url, short_code, title, clicks, is_active, expires_at, created_at, updated_at
		FROM links 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var links []*models.Link
	for rows.Next() {
		link := &models.Link{}
		err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.OriginalURL,
			&link.ShortCode,
			&link.Title,
			&link.Clicks,
			&link.IsActive,
			&link.ExpiresAt,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	
	return links, nil
}

func (r *LinkRepository) Update(link *models.Link) error {
	query := `
		UPDATE links 
		SET original_url = $3, short_code = $4, title = $5, is_active = $6, expires_at = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
	`
	
	return r.db.DB.QueryRow(
		query,
		link.ID,
		link.UserID,
		link.OriginalURL,
		link.ShortCode,
		link.Title,
		link.IsActive,
		link.ExpiresAt,
	).Scan(&link.UpdatedAt)
}

func (r *LinkRepository) Delete(id, userID uuid.UUID) error {
	query := `DELETE FROM links WHERE id = $1 AND user_id = $2`
	result, err := r.db.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("link not found")
	}
	
	return nil
}

func (r *LinkRepository) IncrementClicks(id uuid.UUID) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE id = $1`
	_, err := r.db.DB.Exec(query, id)
	return err
}

func (r *LinkRepository) ShortCodeExists(shortCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM links WHERE short_code = $1)`
	
	err := r.db.DB.QueryRow(query, shortCode).Scan(&exists)
	return exists, err
}

func (r *LinkRepository) GetStats(userID uuid.UUID) (*models.LinkStats, error) {
	stats := &models.LinkStats{}
	
	// Total links
	query := `SELECT COUNT(*) FROM links WHERE user_id = $1`
	err := r.db.DB.QueryRow(query, userID).Scan(&stats.TotalLinks)
	if err != nil {
		return nil, err
	}
	
	// Total clicks
	query = `SELECT COALESCE(SUM(clicks), 0) FROM links WHERE user_id = $1`
	err = r.db.DB.QueryRow(query, userID).Scan(&stats.TotalClicks)
	if err != nil {
		return nil, err
	}
	
	// Active links
	query = `SELECT COUNT(*) FROM links WHERE user_id = $1 AND is_active = true`
	err = r.db.DB.QueryRow(query, userID).Scan(&stats.ActiveLinks)
	if err != nil {
		return nil, err
	}
	
	// Expired links
	query = `SELECT COUNT(*) FROM links WHERE user_id = $1 AND expires_at IS NOT NULL AND expires_at < NOW()`
	err = r.db.DB.QueryRow(query, userID).Scan(&stats.ExpiredLinks)
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}
