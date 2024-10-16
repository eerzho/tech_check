package services

import (
	"context"
	"fmt"
	"tech_check/internal/dto"
	"tech_check/internal/models"
)

type Permission struct {
	permissionRepository PermissionRepository
}

func NewPermission(
	permissionRepository PermissionRepository,
) *Permission {
	return &Permission{
		permissionRepository: permissionRepository,
	}
}

func (p *Permission) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Permission, *dto.Pagination, error) {
	const op = "services.Permission.List"

	permission, pagination, err := p.permissionRepository.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return permission, pagination, nil
}

func (p *Permission) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	const op = "services.Permission.GetByID"

	permission, err := p.permissionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return permission, nil
}
