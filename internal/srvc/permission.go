package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/dto"
	"tech_check/internal/model"
)

type Permission struct {
	permissionRepo PermissionRepo
}

func NewPermission(
	permissionRepo PermissionRepo,
) *Permission {
	return &Permission{
		permissionRepo: permissionRepo,
	}
}

func (p *Permission) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Permission, *dto.Pagination, error) {
	const op = "srvc.Permission.List"

	permission, pagination, err := p.permissionRepo.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return permission, pagination, nil
}

func (p *Permission) GetByID(ctx context.Context, id string) (*model.Permission, error) {
	const op = "srvc.Permission.GetByID"

	permission, err := p.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return permission, nil
}
