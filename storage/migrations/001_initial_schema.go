package migrations

import (
	"context"
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

func CreateInitialSchema(ctx context.Context, driver *ydb.Driver) error {
	tables := []string{
		createUsersTable,
		createArticlesTable,
	}

	for _, tableSQL := range tables {
		err := driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
			return s.ExecuteSchemeQuery(ctx, tableSQL)
		})
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	indexes := []string{createArticleAuthorIndex}

	for _, indexSQL := range indexes {
		err := driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
			return s.ExecuteSchemeQuery(ctx, indexSQL)
		})
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// Table definitions
const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id Int64 NOT NULL,
    notifications_enabled Bool DEFAULT true,
    created_at Timestamp NOT NULL,
    updated_at Timestamp NOT NULL,
    PRIMARY KEY (id)
);
`

const createArticlesTable = `
CREATE TABLE IF NOT EXISTS articles (
    id Utf8 NOT NULL,
    title Utf8 NOT NULL,
    user_id Int64 NOT NULL,
    is_read Bool DEFAULT false,
    created_at Timestamp NOT NULL,
    updated_at Timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    PRIMARY KEY (id)
);
`

// Index definitions
const createArticleAuthorIndex = `
ALTER TABLE articles ADD INDEX user_idx GLOBAL ON (user_id);
`
