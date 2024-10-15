package services

import (
	"context"
	"tech_check/internal/models"
	"time"
)

type (
	UserSrvc interface {
		GetByID(ctx context.Context, id string) (*models.User, error)
		GetByEmail(ctx context.Context, email string) (*models.User, error)
		GetOrCreate(ctx context.Context, email, name, avatar string) (*models.User, error)
	}

	RefreshTokenSrvc interface {
		Delete(ctx context.Context, user *models.User) error
		GetByID(ctx context.Context, user *models.User, id string) (*models.RefreshToken, error)
		Create(ctx context.Context, user *models.User, ip, hash string, expiresAt time.Time) (*models.RefreshToken, error)
	}

	RoleSrvc interface {
		GetByID(ctx context.Context, id string) (*models.Role, error)
	}

	PermissionSrvc interface {
		GetByID(ctx context.Context, id string) (*models.Permission, error)
	}

	CategorySrvc interface {
		GetByID(ctx context.Context, id string) (*models.Category, error)
	}

	QuestionSrvc interface {
		GetRandom(ctx context.Context, category *models.Category, grade string, count int) ([]models.Question, error)
	}

	SessionQuestionSrvc interface {
		Create(ctx context.Context, session *models.Session, text string) (*models.SessionQuestion, error)
	}
)
