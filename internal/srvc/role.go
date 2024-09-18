package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/dto"
	"tech_check/internal/model"

	"github.com/gosimple/slug"
)

type Role struct {
	roleRepo       RoleRepo
	permissionSrvc PermissionSrvc
}

func NewRole(
	roleRepo RoleRepo,
	permissionSrvc PermissionSrvc,
) *Role {
	return &Role{
		roleRepo:       roleRepo,
		permissionSrvc: permissionSrvc,
	}
}

func (r *Role) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Role, *dto.Pagination, error) {
	const op = "srvc.Role.List"

	role, pagination, err := r.roleRepo.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, pagination, nil
}

func (r *Role) Create(ctx context.Context, name string) (*model.Role, error) {
	const op = "srvc.Role.Create"

	slug := slug.Make(name)
	count, err := r.roleRepo.CountBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count+1)
	}

	role := model.Role{
		Name: name,
		Slug: slug,
	}

	err = r.roleRepo.Create(ctx, &role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

func (r *Role) GetByID(ctx context.Context, id string) (*model.Role, error) {
	const op = "srvc.Role.GetByID"

	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Role) Delete(ctx context.Context, id string) error {
	const op = "srvc.Role.Delete"

	err := r.roleRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Role) AddPermission(ctx context.Context, id, permissionID string) (*model.Role, error) {
	const op = "srvc.Role.AddPermission"

	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	permission, err := r.permissionSrvc.GetByID(ctx, permissionID)
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
	err = r.roleRepo.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Role) RemovePermission(ctx context.Context, id, permissionID string) (*model.Role, error) {
	const op = "srvc.Role.RemovePermission"

	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	permission, err := r.permissionSrvc.GetByID(ctx, permissionID)
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
	err = r.roleRepo.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Role) Update(ctx context.Context, id, name string) (*model.Role, error) {
	const op = "srvc.Role.Update"

	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	role.Name = name
	err = r.roleRepo.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
