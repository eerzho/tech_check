package services

import (
	"context"
	"fmt"
	"tech_check/internal/models"
	"time"
)

type RefreshToken struct {
	refreshTokenRepository RefreshTokenRepository
}

func NewRefreshToken(
	refreshTokenRepository RefreshTokenRepository,
) *RefreshToken {
	return &RefreshToken{
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (r *RefreshToken) Delete(ctx context.Context, user *models.User) error {
	const op = "services.RefreshToken.Delete"

	err := r.refreshTokenRepository.Delete(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RefreshToken) Create(ctx context.Context, user *models.User, ip, hash string, expiresAt time.Time) (*models.RefreshToken, error) {
	const op = "services.RefreshToken.Create"

	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		IP:        ip,
		Hash:      hash,
		ExpiresAt: expiresAt,
	}
	err := r.refreshTokenRepository.Create(ctx, &refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &refreshToken, nil
}

func (r *RefreshToken) GetByID(ctx context.Context, user *models.User, id string) (*models.RefreshToken, error) {
	const op = "services.RefreshToken.GetByID"

	refreshToken, err := r.refreshTokenRepository.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}
