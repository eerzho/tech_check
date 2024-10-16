package services

import (
	"context"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"
)

type (
	UserRepository interface {
		Delete(ctx context.Context, id string) error
		Create(ctx context.Context, user *models.User) error
		Update(ctx context.Context, user *models.User) error
		GetByID(ctx context.Context, id string) (*models.User, error)
		ExistsByEmail(ctx context.Context, email string) (bool, error)
		GetByEmail(ctx context.Context, email string) (*models.User, error)
		HasPermission(ctx context.Context, user *models.User, permissionSlug string) (bool, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.User, *dto.Pagination, error)
	}

	RefreshTokenRepository interface {
		Delete(ctx context.Context, user *models.User) error
		Create(ctx context.Context, refreshToken *models.RefreshToken) error
		GetByID(ctx context.Context, user *models.User, id string) (*models.RefreshToken, error)
	}

	RoleRepository interface {
		Delete(ctx context.Context, id string) error
		Create(ctx context.Context, role *models.Role) error
		Update(ctx context.Context, role *models.Role) error
		CountBySlug(ctx context.Context, slug string) (int, error)
		GetByID(ctx context.Context, id string) (*models.Role, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Role, *dto.Pagination, error)
	}

	PermissionRepository interface {
		GetByID(ctx context.Context, id string) (*models.Permission, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Permission, *dto.Pagination, error)
	}

	CategoryRepository interface {
		Delete(ctx context.Context, id string) error
		CountBySlug(ctx context.Context, slug string) (int, error)
		Create(ctx context.Context, category *models.Category) error
		Update(ctx context.Context, category *models.Category) error
		GetByID(ctx context.Context, id string) (*models.Category, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Category, *dto.Pagination, error)
	}

	QuestionRepository interface {
		Delete(ctx context.Context, id string) error
		Create(ctx context.Context, question *models.Question) error
		Update(ctx context.Context, question *models.Question) error
		GetByID(ctx context.Context, id string) (*models.Question, error)
		GetRandom(ctx context.Context, category *models.Category, grade constants.GradeName, count int) ([]models.Question, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Question, *dto.Pagination, error)
	}

	SessionRepository interface {
		Create(ctx context.Context, session *models.Session) error
		Update(ctx context.Context, session *models.Session) error
		GetByID(ctx context.Context, id string) (*models.Session, error)
		ExistsActive(ctx context.Context, user *models.User) (bool, error)
		List(ctx context.Context, user *models.User, page, count int) ([]models.Session, *dto.Pagination, error)
	}

	SessionQuestionRepository interface {
		Create(ctx context.Context, question *models.SessionQuestion) error
		Update(ctx context.Context, question *models.SessionQuestion) error
		List(ctx context.Context, session *models.Session) ([]models.SessionQuestion, error)
		GetByID(ctx context.Context, session *models.Session, id string) (*models.SessionQuestion, error)
	}
)
