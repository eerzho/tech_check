package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/dto"
	"tech_check/internal/model"

	"github.com/gosimple/slug"
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

func (p *Permission) Create(ctx context.Context, name string) (*model.Permission, error) {
	const op = "srvc.Permission.Create"

	slug := slug.Make(name)
	count, err := p.permissionRepo.CountBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count+1)
	}

	permission := model.Permission{
		Name: name,
		Slug: slug,
	}

	err = p.permissionRepo.Create(ctx, &permission)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &permission, nil
}

func (p *Permission) GetByID(ctx context.Context, id string) (*model.Permission, error) {
	const op = "srvc.Permission.GetByID"

	permission, err := p.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return permission, nil
}

func (p *Permission) Delete(ctx context.Context, id string) error {
	const op = "srvc.Permission.Delete"

	err := p.permissionRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Permission) Update(ctx context.Context, id, name string) (*model.Permission, error) {
	const op = "srvc.Permission.Update"

	permission, err := p.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	permission.Name = name
	err = p.permissionRepo.Update(ctx, permission)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return permission, nil
}
