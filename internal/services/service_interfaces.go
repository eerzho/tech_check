package services

import (
	"context"
	"tech_check/internal/models"
	"time"
)

type (
	UserService interface {
		GetByID(ctx context.Context, id string) (*models.User, error)
		GetByEmail(ctx context.Context, email string) (*models.User, error)
		GetOrCreate(ctx context.Context, email, name, avatar string) (*models.User, error)
	}

	RefreshTokenService interface {
		Delete(ctx context.Context, user *models.User) error
		GetByID(ctx context.Context, user *models.User, id string) (*models.RefreshToken, error)
		Create(ctx context.Context, user *models.User, ip, hash string, expiresAt time.Time) (*models.RefreshToken, error)
	}

	RoleService interface {
		GetByID(ctx context.Context, id string) (*models.Role, error)
	}

	PermissionService interface {
		GetByID(ctx context.Context, id string) (*models.Permission, error)
	}

	CategoryService interface {
		GetByID(ctx context.Context, id string) (*models.Category, error)
	}

	QuestionService interface {
		GetRandom(ctx context.Context, category *models.Category, grade string, count int) ([]models.Question, error)
	}

	SessionQuestionService interface {
		Create(ctx context.Context, session *models.Session, text string) (*models.SessionQuestion, error)
	}
)
