package middlewares

import (
	"context"
	"tech_check/internal/dto"
	"tech_check/internal/models"
)

type (
	AuthSrvc interface {
		DecodeAToken(ctx context.Context, aToken string) (*dto.Claims, error)
	}

	UserSrvc interface {
		GetByID(ctx context.Context, id string) (*models.User, error)
		HasPermission(ctx context.Context, user *models.User, permissionSlug string) (bool, error)
	}
)
