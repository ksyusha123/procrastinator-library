package storage

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (s *SQLiteDb) Close() error {
	return s.db.Close()
}

func (s *SQLiteDb) SaveArticle(article *Article) error {
	_, err := s.db.Exec(
		`INSERT INTO articles (url, title, user_id) 
		VALUES (?, ?, ?)`,
		article.URL,
		article.Title,
		article.UserID,
	)
	return err
}

func (s *SQLiteDb) GetArticles(userID int64) ([]Article, error) {
	rows, err := s.db.Query(
		`SELECT id, url, title, is_read, created_at, user_id 
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
		err := rows.Scan(
			&a.ID,
			&a.URL,
			&a.Title,
			&a.IsRead,
			&a.CreatedAt,
			&a.UserID,
		)
		if err != nil {
			log.Printf("Error scanning article row: %v", err)
			continue
		}

		articles = append(articles, a)
	}

	return articles, nil
}

func (s *SQLiteDb) MarkAsRead(articleID int, userID int64) error {
	_, err := s.db.Exec(
		`UPDATE user_articles 
		SET is_read = TRUE 
		WHERE article_id = ? AND user_id = ?`,
		articleID,
		userID,
	)
	return err
}

func (s *SQLiteDb) DeleteArticle(articleID int, userID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM articles 
		WHERE id = ? AND user_id = ?`,
		articleID,
		userID,
	)
	return err
}

func (s *SQLiteDb) GetUnreadArticles(userID int64) ([]Article, error) {
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

func (s *SQLiteDb) GetAllUnreadArticles() ([]Article, error) {
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
