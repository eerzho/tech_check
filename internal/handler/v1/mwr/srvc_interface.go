package mwr

import (
	"context"
	"tech_check/internal/dto"
	"tech_check/internal/model"
)

type (
	AuthSrvc interface {
		DecodeAToken(ctx context.Context, aToken string) (*dto.Claims, error)
	}

	UserSrvc interface {
		GetByID(ctx context.Context, id string) (*model.User, error)
	}
)
