package users

import (
	"context"
)

type UserStorage interface {
	Save(ctx context.Context, userID int64) error
	GetForNotifications(ctx context.Context) ([]int64, error)
}
