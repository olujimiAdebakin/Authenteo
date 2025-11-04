package repository

import (
	"authentio/internal/models"
	"context"
)

// TokenRepository defines the interface for token-related database operations
type TokenRepository interface {
	// SaveRefreshToken stores a new refresh token
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error

	// GetRefreshToken retrieves a refresh token by its token string
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)

	// DeleteRefreshToken removes a refresh token (used during logout or token rotation)
	DeleteRefreshToken(ctx context.Context, token string) error

	// DeleteUserRefreshTokens removes all refresh tokens for a specific user
	DeleteUserRefreshTokens(ctx context.Context, userID int64) error

	// CleanupExpiredTokens removes all expired refresh tokens
	CleanupExpiredTokens(ctx context.Context) error
}