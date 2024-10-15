package srvc

import (
	"context"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"
)

type (
	UserRepo interface {
		Delete(ctx context.Context, id string) error
		Create(ctx context.Context, user *model.User) error
		Update(ctx context.Context, user *model.User) error
		GetByID(ctx context.Context, id string) (*model.User, error)
		ExistsByEmail(ctx context.Context, email string) (bool, error)
		GetByEmail(ctx context.Context, email string) (*model.User, error)
		HasPermission(ctx context.Context, user *model.User, permissionSlug string) (bool, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.User, *dto.Pagination, error)
	}

	RefreshTokenRepo interface {
		Delete(ctx context.Context, user *model.User) error
		Create(ctx context.Context, refreshToken *model.RefreshToken) error
		GetByID(ctx context.Context, user *model.User, id string) (*model.RefreshToken, error)
	}

	RoleRepo interface {
		Delete(ctx context.Context, id string) error
		Create(ctx context.Context, role *model.Role) error
		Update(ctx context.Context, role *model.Role) error
		CountBySlug(ctx context.Context, slug string) (int, error)
		GetByID(ctx context.Context, id string) (*model.Role, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Role, *dto.Pagination, error)
	}

	PermissionRepo interface {
		GetByID(ctx context.Context, id string) (*model.Permission, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Permission, *dto.Pagination, error)
	}

	CategoryRepo interface {
		Delete(ctx context.Context, id string) error
		CountBySlug(ctx context.Context, slug string) (int, error)
		Create(ctx context.Context, category *model.Category) error
		Update(ctx context.Context, category *model.Category) error
		GetByID(ctx context.Context, id string) (*model.Category, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Category, *dto.Pagination, error)
	}

	QuestionRepo interface {
		Delete(ctx context.Context, id string) error
		Create(ctx context.Context, question *model.Question) error
		Update(ctx context.Context, question *model.Question) error
		GetByID(ctx context.Context, id string) (*model.Question, error)
		GetRandom(ctx context.Context, category *model.Category, grade def.GradeName, count int) ([]model.Question, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Question, *dto.Pagination, error)
	}

	SessionRepo interface {
		Create(ctx context.Context, session *model.Session) error
		Update(ctx context.Context, session *model.Session) error
		GetByID(ctx context.Context, id string) (*model.Session, error)
		ExistsActive(ctx context.Context, user *model.User) (bool, error)
		List(ctx context.Context, user *model.User, page, count int) ([]model.Session, *dto.Pagination, error)
	}

	SessionQuestionRepo interface {
		Create(ctx context.Context, question *model.SessionQuestion) error
		Update(ctx context.Context, question *model.SessionQuestion) error
		List(ctx context.Context, session *model.Session) ([]model.SessionQuestion, error)
		GetByID(ctx context.Context, session *model.Session, id string) (*model.SessionQuestion, error)
	}
)
