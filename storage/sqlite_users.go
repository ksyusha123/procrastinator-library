package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type UserSQLiteStorage struct {
	db *sql.DB
}

func (s *SQLiteDb) SaveUser(userID int64) error {
	_, err := s.db.Exec(
		`INSERT INTO users (id) 
		VALUES (?)`,
		userID,
	)
	return err
}

func (s *SQLiteDb) GetUsersReceivingNotifications() ([]User, error) {
	rows, err := s.db.Query(
		`SELECT id 
		FROM users 
		WHERE notifications_enabled = TRUE`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
		)
		if err != nil {
			log.Printf("Error scanning article row: %v", err)
			continue
		}

		users = append(users, u)
	}

	return users, nil
}
