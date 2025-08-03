package articles

import (
	"context"
	"github.com/google/uuid"
)

type ArticleStorage interface {
	Save(ctx context.Context, article *Article) error
	Get(ctx context.Context, userID int64) ([]Article, error)
	MarkAsRead(ctx context.Context, articleID uuid.UUID, userID int64) error
	Delete(ctx context.Context, articleID uuid.UUID, userID int64) error
	GetUnread(ctx context.Context, userID int64) ([]Article, error)
}
