package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/model"
	"time"
)

type RefreshToken struct {
	refreshTokenRepo RefreshTokenRepo
}

func NewRefreshToken(
	refreshTokenRepo RefreshTokenRepo,
) *RefreshToken {
	return &RefreshToken{
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (r *RefreshToken) Delete(ctx context.Context, user *model.User) error {
	const op = "srvc.RefreshToken.Delete"

	err := r.refreshTokenRepo.Delete(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RefreshToken) Create(ctx context.Context, user *model.User, ip, hash string, expiresAt time.Time) (*model.RefreshToken, error) {
	const op = "srvc.RefreshToken.Create"

	refreshToken := model.RefreshToken{
		UserID:    user.ID,
		IP:        ip,
		Hash:      hash,
		ExpiresAt: expiresAt,
	}
	err := r.refreshTokenRepo.Create(ctx, &refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &refreshToken, nil
}

func (r *RefreshToken) GetByID(ctx context.Context, user *model.User, id string) (*model.RefreshToken, error) {
	const op = "srvc.RefreshToken.GetByID"

	refreshToken, err := r.refreshTokenRepo.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}
