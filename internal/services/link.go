package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"link-shortener/internal/models"
	"link-shortener/internal/repository"
	"link-shortener/internal/utils"
)

type LinkService struct {
	linkRepo *repository.LinkRepository
	baseURL  string
}

func NewLinkService(linkRepo *repository.LinkRepository, baseURL string) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
		baseURL:  strings.TrimSuffix(baseURL, "/"),
	}
}

func (s *LinkService) CreateLink(userID uuid.UUID, req *models.CreateLinkRequest) (*models.LinkResponse, error) {
	// Validate and sanitize URL
	if err := utils.ValidateURL(req.OriginalURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	
	originalURL := utils.SanitizeURL(req.OriginalURL)

	// Generate or use custom short code
	var shortCode string
	if req.CustomAlias != "" {
		// Validate custom alias
		if err := utils.ValidateShortCode(req.CustomAlias); err != nil {
			return nil, fmt.Errorf("invalid custom alias: %w", err)
		}
		
		// Check if custom alias already exists
		exists, err := s.linkRepo.ShortCodeExists(req.CustomAlias)
		if err != nil {
			return nil, fmt.Errorf("failed to check short code: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("custom alias already exists")
		}
		shortCode = req.CustomAlias
	} else {
		// Generate random short code
		for {
			generatedCode, err := utils.GenerateShortCode(8)
			if err != nil {
				return nil, fmt.Errorf("failed to generate short code: %w", err)
			}
			
			exists, err := s.linkRepo.ShortCodeExists(generatedCode)
			if err != nil {
				return nil, fmt.Errorf("failed to check short code: %w", err)
			}
			
			if !exists {
				shortCode = generatedCode
				break
			}
		}
	}

	// Create link
	link := &models.Link{
		ID:          uuid.New(),
		UserID:      userID,
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		Title:       req.Title,
		ExpiresAt:   req.ExpiresAt,
		IsActive:    true,
	}

	if err := s.linkRepo.Create(link); err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	return s.toLinkResponse(link), nil
}

func (s *LinkService) GetLinkByID(userID, linkID uuid.UUID) (*models.LinkResponse, error) {
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return nil, fmt.Errorf("link not found: %w", err)
	}

	// Check if user owns this link
	if link.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	return s.toLinkResponse(link), nil
}

func (s *LinkService) GetLinksByUserID(userID uuid.UUID, limit, offset int) ([]*models.LinkResponse, error) {
	links, err := s.linkRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get links: %w", err)
	}

	var responses []*models.LinkResponse
	for _, link := range links {
		responses = append(responses, s.toLinkResponse(link))
	}

	return responses, nil
}

func (s *LinkService) UpdateLink(userID, linkID uuid.UUID, req *models.UpdateLinkRequest) (*models.LinkResponse, error) {
	// Get existing link
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return nil, fmt.Errorf("link not found: %w", err)
	}

	// Check if user owns this link
	if link.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update fields if provided
	if req.OriginalURL != "" {
		if err := utils.ValidateURL(req.OriginalURL); err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		link.OriginalURL = utils.SanitizeURL(req.OriginalURL)
	}

	if req.CustomAlias != "" {
		if err := utils.ValidateShortCode(req.CustomAlias); err != nil {
			return nil, fmt.Errorf("invalid custom alias: %w", err)
		}
		
		// Check if new alias already exists (excluding current link)
		exists, err := s.linkRepo.ShortCodeExists(req.CustomAlias)
		if err != nil {
			return nil, fmt.Errorf("failed to check short code: %w", err)
		}
		if exists && req.CustomAlias != link.ShortCode {
			return nil, fmt.Errorf("custom alias already exists")
		}
		link.ShortCode = req.CustomAlias
	}

	if req.Title != "" {
		link.Title = req.Title
	}

	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}

	if req.ExpiresAt != nil {
		link.ExpiresAt = req.ExpiresAt
	}

	// Update link
	if err := s.linkRepo.Update(link); err != nil {
		return nil, fmt.Errorf("failed to update link: %w", err)
	}

	return s.toLinkResponse(link), nil
}

func (s *LinkService) DeleteLink(userID, linkID uuid.UUID) error {
	// Check if link exists and user owns it
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return fmt.Errorf("link not found: %w", err)
	}

	if link.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.linkRepo.Delete(linkID, userID)
}

func (s *LinkService) RedirectToOriginal(shortCode string) (string, error) {
	link, err := s.linkRepo.GetByShortCode(shortCode)
	if err != nil {
		return "", fmt.Errorf("link not found: %w", err)
	}

	// Increment click count
	go func() {
		if err := s.linkRepo.IncrementClicks(link.ID); err != nil {
			// Log error but don't fail the redirect
			fmt.Printf("Failed to increment clicks: %v\n", err)
		}
	}()

	return link.OriginalURL, nil
}

func (s *LinkService) GetStats(userID uuid.UUID) (*models.LinkStats, error) {
	return s.linkRepo.GetStats(userID)
}

func (s *LinkService) toLinkResponse(link *models.Link) *models.LinkResponse {
	return &models.LinkResponse{
		ID:          link.ID,
		OriginalURL: link.OriginalURL,
		ShortCode:   link.ShortCode,
		ShortURL:    fmt.Sprintf("%s/r/%s", s.baseURL, link.ShortCode),
		Title:       link.Title,
		Clicks:      link.Clicks,
		IsActive:    link.IsActive,
		ExpiresAt:   link.ExpiresAt,
		CreatedAt:   link.CreatedAt,
		UpdatedAt:   link.UpdatedAt,
	}
}
