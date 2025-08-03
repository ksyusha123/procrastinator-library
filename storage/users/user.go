package users

import "time"

type User struct {
	ID                   int64
	NotificationsEnabled bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
