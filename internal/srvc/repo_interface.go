package srvc

import (
	"context"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"
)

type (
	UserRepo interface {
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.User, *dto.Pagination, error)
		Create(ctx context.Context, user *model.User) error
		GetByID(ctx context.Context, id string) (*model.User, error)
		Update(ctx context.Context, user *model.User) error
		Delete(ctx context.Context, id string) error
		IsExistsEmail(ctx context.Context, email string) (bool, error)
		GetByEmail(ctx context.Context, email string) (*model.User, error)
		HasPermission(ctx context.Context, user *model.User, permissionSlug string) (bool, error)
	}

	RefreshTokenRepo interface {
		Create(ctx context.Context, refreshToken *model.RefreshToken) error
		GetByUserAndID(ctx context.Context, user *model.User, id string) (*model.RefreshToken, error)
		DeleteByUser(ctx context.Context, user *model.User) error
	}

	RoleRepo interface {
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Role, *dto.Pagination, error)
		Create(ctx context.Context, role *model.Role) error
		GetByID(ctx context.Context, id string) (*model.Role, error)
		Update(ctx context.Context, role *model.Role) error
		Delete(ctx context.Context, id string) error
		CountBySlug(ctx context.Context, slug string) (int, error)
	}

	PermissionRepo interface {
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Permission, *dto.Pagination, error)
		GetByID(ctx context.Context, id string) (*model.Permission, error)
	}

	CategoryRepo interface {
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Category, *dto.Pagination, error)
		Create(ctx context.Context, category *model.Category) error
		GetByID(ctx context.Context, id string) (*model.Category, error)
		Update(ctx context.Context, category *model.Category) error
		Delete(ctx context.Context, id string) error
		CountBySlug(ctx context.Context, slug string) (int, error)
	}

	QuestionRepo interface {
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Question, *dto.Pagination, error)
		Create(ctx context.Context, question *model.Question) error
		GetByID(ctx context.Context, id string) (*model.Question, error)
		Update(ctx context.Context, question *model.Question) error
		Delete(ctx context.Context, id string) error
		GetRandom(ctx context.Context, category *model.Category, grade def.GradeName, count int) ([]model.Question, error)
	}

	SessionRepo interface {
		List(ctx context.Context, user *model.User, page, count int) ([]model.Session, *dto.Pagination, error)
		Create(ctx context.Context, session *model.Session) error
		GetByID(ctx context.Context, id string) (*model.Session, error)
		Update(ctx context.Context, session *model.Session) error
		IsExistsActive(ctx context.Context, user *model.User) (bool, error)
	}

	SessionQuestionRepo interface {
		List(ctx context.Context, session *model.Session) ([]model.SessionQuestion, error)
		Create(ctx context.Context, question *model.SessionQuestion) error
		GetByID(ctx context.Context, session *model.Session, id string) (*model.SessionQuestion, error)
		Update(ctx context.Context, question *model.SessionQuestion) error
	}
)
