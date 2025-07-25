package storage

type User struct {
	ID                   int64
	NotificationsEnabled bool
}

type UserStorage interface {
	SaveUser(userID int64) error
	GetUsersReceivingNotifications() ([]User, error)
}
