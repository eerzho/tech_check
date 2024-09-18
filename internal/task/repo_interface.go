package task

import "context"

type (
	UserRepo interface {
		RemoveRole(ctx context.Context, roleID string) error
	}
)
