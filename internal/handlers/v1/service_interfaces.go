package v1

import (
	"context"
	"tech_check/internal/dto"
	"tech_check/internal/models"
)

type (
	UserService interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*models.User, error)
		Update(ctx context.Context, id, name string) (*models.User, error)
		AddRole(ctx context.Context, id, roleID string) (*models.User, error)
		RemoveRole(ctx context.Context, id, roleID string) (*models.User, error)
		Create(ctx context.Context, email, name, password string) (*models.User, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.User, *dto.Pagination, error)
	}

	AuthService interface {
		DecodeAToken(ctx context.Context, aToken string) (*dto.Claims, error)
		GoogleLogin(ctx context.Context, tokenID, ip string) (*dto.Token, error)
		Login(ctx context.Context, email, password, ip string) (*dto.Token, error)
		Refresh(ctx context.Context, aToken, rToken, ip string) (*dto.Token, error)
	}

	RoleService interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*models.Role, error)
		Create(ctx context.Context, name string) (*models.Role, error)
		Update(ctx context.Context, id, name string) (*models.Role, error)
		AddPermission(ctx context.Context, id, permissionID string) (*models.Role, error)
		RemovePermission(ctx context.Context, id, permissionID string) (*models.Role, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Role, *dto.Pagination, error)
	}

	PermissionService interface {
		GetByID(ctx context.Context, id string) (*models.Permission, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Permission, *dto.Pagination, error)
	}

	CategoryService interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*models.Category, error)
		Create(ctx context.Context, name, description string) (*models.Category, error)
		Update(ctx context.Context, id, name, description string) (*models.Category, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Category, *dto.Pagination, error)
	}

	QuestionService interface {
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*models.Question, error)
		Update(ctx context.Context, id, grade, text string) (*models.Question, error)
		Create(ctx context.Context, text, grade, categoryID string) (*models.Question, error)
		List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Question, *dto.Pagination, error)
	}

	SessionService interface {
		Cancel(ctx context.Context, user *models.User, id string) (*models.Session, error)
		GetByID(ctx context.Context, user *models.User, id string) (*models.Session, error)
		Summarize(ctx context.Context, user *models.User, id string) (*models.Session, error)
		Create(ctx context.Context, user *models.User, categoryID, grade string) (*models.Session, error)
		List(ctx context.Context, user *models.User, page, count int) ([]models.Session, *dto.Pagination, error)
	}

	SessionQuestionService interface {
		List(ctx context.Context, session *models.Session) ([]models.SessionQuestion, error)
		GetByID(ctx context.Context, session *models.Session, id string) (*models.SessionQuestion, error)
		Update(ctx context.Context, session *models.Session, id, answer string) (*models.SessionQuestion, error)
	}
)
