package v1

import (
	"context"
	"tech_check/internal/dto"
	"tech_check/internal/model"
)

type (
	UserSrvc interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*model.User, error)
		Update(ctx context.Context, id, name string) (*model.User, error)
		AddRole(ctx context.Context, id, roleID string) (*model.User, error)
		RemoveRole(ctx context.Context, id, roleID string) (*model.User, error)
		Create(ctx context.Context, email, name, password string) (*model.User, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.User, *dto.Pagination, error)
	}

	AuthSrvc interface {
		DecodeAToken(ctx context.Context, aToken string) (*dto.Claims, error)
		GoogleLogin(ctx context.Context, tokenID, ip string) (*dto.Token, error)
		Login(ctx context.Context, email, password, ip string) (*dto.Token, error)
		Refresh(ctx context.Context, aToken, rToken, ip string) (*dto.Token, error)
	}

	RoleSrvc interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*model.Role, error)
		Create(ctx context.Context, name string) (*model.Role, error)
		Update(ctx context.Context, id, name string) (*model.Role, error)
		AddPermission(ctx context.Context, id, permissionID string) (*model.Role, error)
		RemovePermission(ctx context.Context, id, permissionID string) (*model.Role, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Role, *dto.Pagination, error)
	}

	PermissionSrvc interface {
		GetByID(ctx context.Context, id string) (*model.Permission, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Permission, *dto.Pagination, error)
	}

	CategorySrvc interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*model.Category, error)
		Create(ctx context.Context, name, description string) (*model.Category, error)
		Update(ctx context.Context, id, name, description string) (*model.Category, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Category, *dto.Pagination, error)
	}

	QuestionSrvc interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*model.Question, error)
		Update(ctx context.Context, id, grade, text string) (*model.Question, error)
		Create(ctx context.Context, text, grade, categoryID string) (*model.Question, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Question, *dto.Pagination, error)
	}

	SessionSrvc interface {
		Cancel(ctx context.Context, user *model.User, id string) (*model.Session, error)
		GetByID(ctx context.Context, user *model.User, id string) (*model.Session, error)
		Summarize(ctx context.Context, user *model.User, id string) (*model.Session, error)
		Create(ctx context.Context, user *model.User, categoryID, grade string) (*model.Session, error)
		List(ctx context.Context, user *model.User, page, count int) ([]model.Session, *dto.Pagination, error)
	}

	SessionQuestionSrvc interface {
		List(ctx context.Context, session *model.Session) ([]model.SessionQuestion, error)
		GetByID(ctx context.Context, session *model.Session, id string) (*model.SessionQuestion, error)
		Update(ctx context.Context, session *model.Session, id, answer string) (*model.SessionQuestion, error)
	}
)
