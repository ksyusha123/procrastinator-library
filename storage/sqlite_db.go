package storage

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL UNIQUE,
		title TEXT,
		summary TEXT,
		is_read BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER NOT NULL
	);`)

	return err
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func (s *SQLiteStorage) SaveArticle(article *Article) error {
	_, err := s.db.Exec(
		`INSERT INTO articles (url, title, summary, is_read, user_id) 
		VALUES (?, ?, ?, ?, ?)`,
		article.URL,
		article.Title,
		article.Summary,
		article.IsRead,
		article.UserID,
	)
	return err
}

func (s *SQLiteStorage) GetArticles(userID int64) ([]Article, error) {
	rows, err := s.db.Query(
		`SELECT id, url, title, summary, is_read, created_at, user_id 
		FROM articles 
		WHERE user_id = ? 
		ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		var createdAt string
		err := rows.Scan(
			&a.ID,
			&a.URL,
			&a.Title,
			&a.Summary,
			&a.IsRead,
			&createdAt,
			&a.UserID,
		)
		if err != nil {
			log.Printf("Error scanning article row: %v", err)
			continue
		}

		a.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		articles = append(articles, a)
	}

	return articles, nil
}

func (s *SQLiteStorage) MarkAsRead(articleID int, userID int64) error {
	_, err := s.db.Exec(
		`UPDATE articles 
		SET is_read = TRUE 
		WHERE id = ? AND user_id = ?`,
		articleID,
		userID,
	)
	return err
}

func (s *SQLiteStorage) DeleteArticle(articleID int, userID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM articles 
		WHERE id = ? AND user_id = ?`,
		articleID,
		userID,
	)
	return err
}

func (s *SQLiteStorage) GetUnreadArticles(userID int64) ([]Article, error) {
	rows, err := s.db.Query(
		`SELECT id, url, title 
		FROM articles 
		WHERE user_id = ? AND is_read = FALSE`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.URL, &a.Title)
		if err != nil {
			log.Printf("Error scanning unread article: %v", err)
			continue
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (s *SQLiteStorage) GetAllUnreadArticles() ([]Article, error) {
	rows, err := s.db.Query(
		`SELECT id, url, title, user_id 
		FROM articles 
		WHERE is_read = FALSE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.URL, &a.Title, &a.UserID)
		if err != nil {
			log.Printf("Error scanning unread article: %v", err)
			continue
		}
		articles = append(articles, a)
	}

	return articles, nil
}

// func (s *SQLiteStorage) GetUsers() ([]User, error) {
// 	rows, err := s.db.Query(`SELECT id FROM users`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next(){
// 		var u User
// 		err := rows.Scan(&u.ID)
// 		if err != nil {
// 			log.Printf("Error scanning user: %v", err)
// 			continue
// 		}
// 		users = append(users, u)
// 	}
// 	return users, nil
// }