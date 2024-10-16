package services

import (
	"context"
	"errors"
	"fmt"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepository UserRepository
	roleService    RoleService
}

func NewUser(
	userRepository UserRepository,
	roleService RoleService,
) *User {
	return &User{
		userRepository: userRepository,
		roleService:    roleService,
	}
}

func (u *User) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.User, *dto.Pagination, error) {
	const op = "services.User.List"

	user, pagination, err := u.userRepository.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, pagination, nil
}

func (u *User) Create(ctx context.Context, email, name, password string) (*models.User, error) {
	const op = "services.User.Create"

	exists, err := u.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrAlreadyExists)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user := models.User{
		Email:    email,
		Name:     name,
		Password: string(passwordHash),
	}

	err = u.userRepository.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (u *User) GetOrCreate(ctx context.Context, email, name, avatar string) (*models.User, error) {
	const op = "services.User.GetOrCreate"

	user, err := u.GetByEmail(ctx, email)
	if err == nil {
		user.Name = name
		user.Avatar = avatar

		err = u.userRepository.Update(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return user, nil
	}

	if !errors.Is(err, constants.ErrNotFound) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user = &models.User{
		Email:  email,
		Name:   name,
		Avatar: avatar,
	}
	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) GetByID(ctx context.Context, id string) (*models.User, error) {
	const op = "services.User.GetByID"

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) Update(ctx context.Context, id, name string) (*models.User, error) {
	const op = "services.User.Update"

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Name = name
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) Delete(ctx context.Context, id string) error {
	const op = "services.User.Delete"

	err := u.userRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "services.User.GetByEmail"

	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) AddRole(ctx context.Context, id, roleID string) (*models.User, error) {
	const op = "services.User.AddRole"

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	role, err := u.roleService.GetByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	isUnique := true
	for _, roleID := range user.RoleIDs {
		if roleID == role.ID {
			isUnique = false
			break
		}
	}
	if !isUnique {
		return user, nil
	}

	user.RoleIDs = append(user.RoleIDs, role.ID)
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) RemoveRole(ctx context.Context, id, roleID string) (*models.User, error) {
	const op = "services.User.RemoveRole"

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	role, err := u.roleService.GetByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	existsIdx := -1
	for index, roleID := range user.RoleIDs {
		if roleID == role.ID {
			existsIdx = index
			break
		}
	}
	if existsIdx == -1 {
		return user, nil
	}

	user.RoleIDs = append(user.RoleIDs[:existsIdx], user.RoleIDs[existsIdx+1:]...)
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) HasPermission(ctx context.Context, user *models.User, permissionSlug string) (bool, error) {
	const op = "services.User.HasPermission"

	has, err := u.userRepository.HasPermission(ctx, user, permissionSlug)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return has, nil
}
