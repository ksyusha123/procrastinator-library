package users

import (
	"context"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"time"
)

type YDBUserStorage struct {
	db *ydb.Driver
}

func (s *YDBUserStorage) Save(ctx context.Context, telegramUserID int64) error {
	query := `
	DECLARE $id AS Int64;
	DECLARE $notifications_enabled AS Bool;
	DECLARE $created_at AS Timestamp;
	
	UPSERT INTO users (id, notifications_enabled, created_at, updated_at)
	VALUES ($id, $notifications_enabled, $created_at, $created_at);
	`

	return s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, _, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$id", types.Int64Value(telegramUserID)),
				table.ValueParam("$notifications_enabled", types.BoolValue(true)),
				table.ValueParam("$created_at", types.TimestampValueFromTime(time.Now().UTC())),
			),
		)
		return err
	})
}

func (s *YDBUserStorage) GetForNotifications(ctx context.Context) ([]int64, error) {
	query := `-- noinspection SqlNoDataSourceInspectionForFile
	
		SELECT id
		FROM users
		WHERE notifications_enabled = true;
		`

	var ids []int64

	err := s.db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, nil)
		if err != nil {
			return err
		}
		defer res.Close()

		for res.NextResultSet(ctx) {
			for res.NextRow() {
				var id int64
				if err := res.Scan(&id); err != nil {
					return err
				}
				ids = append(ids, id)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ids, nil

}

func NewYDBUserStorage(db *ydb.Driver) UserStorage {
	return &YDBUserStorage{db: db}
}
