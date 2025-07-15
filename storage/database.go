package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	ID        int
	URL       string
	Title     string
	Summary   string
	Tags      []string
	IsRead    bool
	CreatedAt time.Time
	UserID    int64
}

type Storage interface {
	SaveArticle(article *Article) error
	GetArticles(userID int64) ([]Article, error)
	MarkAsRead(articleID int, userID int64) error
	DeleteArticle(articleID int, userID int64) error
	// SearchByTag(userID int64, tag string) ([]Article, error)
	GetUnreadArticles(userID int64) ([]Article, error)
}

type SQLiteDB struct {
	db *sql.DB
}