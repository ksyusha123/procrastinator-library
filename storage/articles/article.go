package articles

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	ID        uuid.UUID
	URL       string
	Title     string
	UserID    int64
	IsRead    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
