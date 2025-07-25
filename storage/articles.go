package storage

import (
	"time"
)

type Article struct {
	ID        int
	URL       string
	Title     string
	UserID    int64
	IsRead    bool
	CreatedAt time.Time
}

type ArticleSummary struct {
	ID      int
	Summary string
}

type ArticleTags struct {
	ID   int
	Tags []string
}

type ArticleStorage interface {
	SaveArticle(article *Article) error
	GetArticles(userID int64) ([]Article, error)
	MarkAsRead(articleID int, userID int64) error
	DeleteArticle(articleID int, userID int64) error
	GetUnreadArticles(userID int64) ([]Article, error)
}
