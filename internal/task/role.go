package task

import (
	"context"
	"fmt"
)

type RemoveRole struct {
	roleID   string
	userRepo UserRepo // TODO: ??? может лучше сразу создать в NewRemoveRole и не ебаться с прокидыванием через roleRepo ???
}

func NewRemoveRole(userRepo UserRepo, roleID string) *RemoveRole {
	return &RemoveRole{
		roleID:   roleID,
		userRepo: userRepo,
	}
}

func (r *RemoveRole) Execute(ctx context.Context) error {
	const op = "task.RemoveRole.Execute"

	err := r.userRepo.RemoveRole(ctx, r.roleID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
