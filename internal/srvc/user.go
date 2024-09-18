package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepo UserRepo
	roleSrvc RoleSrvc
}

func NewUser(
	userRepo UserRepo,
	roleSrvc RoleSrvc,
) *User {
	return &User{
		userRepo: userRepo,
		roleSrvc: roleSrvc,
	}
}

func (u *User) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.User, *dto.Pagination, error) {
	const op = "srvc.User.List"

	user, pagination, err := u.userRepo.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, pagination, nil
}

func (u *User) Create(ctx context.Context, email, name, password string) (*model.User, error) {
	const op = "srvc.User.Create"

	exists, err := u.userRepo.IsExistsEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return nil, fmt.Errorf("%s: %w", op, def.ErrAlreadyExists)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user := model.User{
		Email:    email,
		Name:     name,
		Password: string(passwordHash),
	}

	err = u.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (u *User) GetByID(ctx context.Context, id string) (*model.User, error) {
	const op = "srvc.User.GetByID"

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) Update(ctx context.Context, id, name string) (*model.User, error) {
	const op = "srvc.User.Update"

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Name = name
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) Delete(ctx context.Context, id string) error {
	const op = "srvc.User.Delete"

	err := u.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "srvc.User.GetByEmail"

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) AddRole(ctx context.Context, id, roleID string) (*model.User, error) {
	const op = "srvc.User.AddRole"

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	role, err := u.roleSrvc.GetByID(ctx, roleID)
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
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) RemoveRole(ctx context.Context, id, roleID string) (*model.User, error) {
	const op = "srvc.User.RemoveRole"

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	role, err := u.roleSrvc.GetByID(ctx, roleID)
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
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
