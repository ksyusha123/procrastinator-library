package articles

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type YDBArticlesStorage struct {
	db *ydb.Driver
}

func (s *YDBArticlesStorage) GetUnread(ctx context.Context, userID int64) ([]Article, error) {
	query := `
	DECLARE $user_id AS Int64;
	
	SELECT 
		id, url, title, user_id, is_read, created_at, updated_at
	FROM articles
	WHERE user_id = $user_id AND is_read = false;
	ORDER BY created_at DESC;
	`

	var articles []Article

	err := s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, res, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$user_id", types.Int64Value(userID)),
			),
		)
		if err != nil {
			return err
		}
		defer res.Close()

		return scanArticles(ctx, res, &articles)
	})

	return articles, err
}

func NewYDBArticleStorage(db *ydb.Driver) ArticleStorage {
	return &YDBArticlesStorage{db: db}
}

func (s *YDBArticlesStorage) Save(ctx context.Context, article *Article) error {
	query := `
	DECLARE $id AS Utf8;
	DECLARE $url AS Utf8;
	DECLARE $title AS Utf8;
	DECLARE $user_id AS Int64;
	DECLARE $is_read AS Bool;
	DECLARE $created_at AS Timestamp;
	
	UPSERT INTO articles (id, url, title, user_id, is_read, created_at, updated_at)
	VALUES ($id, $url, $title, $user_id, $is_read, $created_at, $created_at);
	`

	now := time.Now().UTC()

	return s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, _, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$id", types.UTF8Value(article.ID.String())),
				table.ValueParam("$url", types.UTF8Value(article.URL)),
				table.ValueParam("$title", types.UTF8Value(article.Title)),
				table.ValueParam("$user_id", types.Int64Value(article.UserID)),
				table.ValueParam("$is_read", types.BoolValue(article.IsRead)),
				table.ValueParam("$created_at", types.TimestampValueFromTime(now)),
			),
		)
		return err
	})
}

func (s *YDBArticlesStorage) Get(ctx context.Context, userID int64) ([]Article, error) {
	query := `
	DECLARE $user_id AS Int64;
	
	SELECT 
		id, url, title, user_id, is_read, created_at, updated_at
	FROM articles
	WHERE user_id = $user_id
	ORDER BY created_at DESC;
	`

	var articles []Article

	err := s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, res, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$user_id", types.Int64Value(userID)),
			),
		)
		if err != nil {
			return err
		}
		defer res.Close()

		return scanArticles(ctx, res, &articles)
	})

	return articles, err
}

func (s *YDBArticlesStorage) MarkAsRead(ctx context.Context, articleID uuid.UUID, userID int64) error {
	query := `
	DECLARE $id AS Utf8;
	DECLARE $user_id AS Int64;
	DECLARE $updated_at AS Timestamp;
	
	UPDATE articles
	SET 
		is_read = true,
		updated_at = CurrentUtcTimestamp();
	WHERE id = $id AND user_id = $user_id;
	`

	return s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, _, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$id", types.UTF8Value(articleID.String())),
				table.ValueParam("$user_id", types.Int64Value(userID)),
				table.ValueParam("$updated_at", types.TimestampValueFromTime(time.Now().UTC())),
			),
		)
		return err
	})
}

func (s *YDBArticlesStorage) Delete(ctx context.Context, articleID uuid.UUID, userID int64) error {
	query := `
	DECLARE $id AS Utf8;
	DECLARE $user_id AS Int64;
	
	DELETE FROM articles
	WHERE id = $id AND user_id = $user_id;
	`

	return s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, _, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$id", types.UTF8Value(articleID.String())),
				table.ValueParam("$user_id", types.Int64Value(userID)),
			),
		)
		return err
	})
}

func scanArticles(ctx context.Context, res result.Result, articles *[]Article) error {
	for res.NextResultSet(ctx) {
		for res.NextRow() {
			var idStr string
			var article Article

			if err := res.Scan(
				&idStr,
				&article.URL,
				&article.Title,
				&article.UserID,
				&article.IsRead,
				&article.CreatedAt,
				&article.UpdatedAt,
			); err != nil {
				return err
			}

			id, err := uuid.Parse(idStr)
			if err != nil {
				return fmt.Errorf("failed to parse UUID: %w", err)
			}

			article.ID = id
			*articles = append(*articles, article)
		}
	}
	return nil
}
