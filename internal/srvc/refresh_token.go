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

func NewRefreshToken(refreshTokenRepo RefreshTokenRepo) *RefreshToken {
	return &RefreshToken{
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (r *RefreshToken) DeleteByUser(ctx context.Context, user *model.User) error {
	const op = "srvc.RefreshToken.DeleteByUser"

	err := r.refreshTokenRepo.DeleteByUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RefreshToken) CreateByUser(ctx context.Context, user *model.User, ip, hash string, expiresAt time.Time) (*model.RefreshToken, error) {
	const op = "srvc.RefreshToken.CreateByUser"

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

func (r *RefreshToken) GetByUserAndID(ctx context.Context, user *model.User, id string) (*model.RefreshToken, error) {
	const op = "srvc.RefreshToken.GetByUserAndID"

	refreshToken, err := r.refreshTokenRepo.GetByUserAndID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}
