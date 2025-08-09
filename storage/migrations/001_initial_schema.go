package migrations

import (
	"context"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"path"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

func CreateInitialSchema(ctx context.Context, db *ydb.Driver) error {
	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			return s.CreateTable(ctx, path.Join(db.Name(), "articles"),
				options.WithColumn("id", types.TypeUTF8),
				options.WithColumn("url", types.TypeUTF8),
				options.WithColumn("title", types.TypeUTF8),
				options.WithColumn("user_id", types.TypeInt64),
				options.WithColumn("is_read", types.TypeBool),
				options.WithColumn("created_at", types.TypeTimestamp),
				options.WithColumn("updated_at", types.TypeTimestamp),
				options.WithPrimaryKeyColumn("id"),
				options.WithIndex("idx_articles_user_id",
					options.WithIndexColumns("user_id")),
			)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create articles table: %w", err)
	}

	err = db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			return s.CreateTable(ctx, path.Join(db.Name(), "users"),
				options.WithColumn("id", types.TypeInt64),
				options.WithColumn("notifications_enabled", types.TypeBool),
				options.WithColumn("created_at", types.TypeTimestamp),
				options.WithColumn("updated_at", types.TypeTimestamp),
				options.WithPrimaryKeyColumn("id"),
			)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	//indexes := []string{createArticleAuthorIndex}

	//for _, indexSQL := range indexes {
	//	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
	//		return s.ExecuteSchemeQuery(ctx, indexSQL)
	//	})
	//	if err != nil {
	//		return fmt.Errorf("failed to create index: %w", err)
	//	}
	//}

	return nil
}

// Index definitions
//const createArticleAuthorIndex = `
//ALTER TABLE articles ADD INDEX user_idx GLOBAL ON (user_id);
//`
