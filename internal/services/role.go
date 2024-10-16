package services

import (
	"context"
	"fmt"
	"tech_check/internal/dto"
	"tech_check/internal/models"

	"github.com/gosimple/slug"
)

type Role struct {
	roleRepository    RoleRepository
	permissionService PermissionService
}

func NewRole(
	roleRepository RoleRepository,
	permissionService PermissionService,
) *Role {
	return &Role{
		roleRepository:    roleRepository,
		permissionService: permissionService,
	}
}

func (r *Role) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Role, *dto.Pagination, error) {
	const op = "services.Role.List"

	role, pagination, err := r.roleRepository.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, pagination, nil
}

func (r *Role) Create(ctx context.Context, name string) (*models.Role, error) {
	const op = "services.Role.Create"

	slug := slug.Make(name)
	count, err := r.roleRepository.CountBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count+1)
	}

	role := models.Role{
		Name: name,
		Slug: slug,
	}

	err = r.roleRepository.Create(ctx, &role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

func (r *Role) GetByID(ctx context.Context, id string) (*models.Role, error) {
	const op = "services.Role.GetByID"

	role, err := r.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Role) Delete(ctx context.Context, id string) error {
	const op = "services.Role.Delete"

	err := r.roleRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Role) AddPermission(ctx context.Context, id, permissionID string) (*models.Role, error) {
	const op = "services.Role.AddPermission"

	role, err := r.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	permission, err := r.permissionService.GetByID(ctx, permissionID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	isUnique := true
	for _, permissionID := range role.PermissionIDs {
		if permissionID == permission.ID {
			isUnique = false
			break
		}
	}
	if !isUnique {
		return role, nil
	}

	role.PermissionIDs = append(role.PermissionIDs, permission.ID)
	err = r.roleRepository.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Role) RemovePermission(ctx context.Context, id, permissionID string) (*models.Role, error) {
	const op = "services.Role.RemovePermission"

	role, err := r.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	permission, err := r.permissionService.GetByID(ctx, permissionID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	existsIdx := -1
	for index, permissionID := range role.PermissionIDs {
		if permissionID == permission.ID {
			existsIdx = index
			break
		}
	}
	if existsIdx == -1 {
		return role, nil
	}

	role.PermissionIDs = append(role.PermissionIDs[:existsIdx], role.PermissionIDs[existsIdx+1:]...)
	err = r.roleRepository.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Role) Update(ctx context.Context, id, name string) (*models.Role, error) {
	const op = "services.Role.Update"

	role, err := r.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	role.Name = name
	err = r.roleRepository.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
