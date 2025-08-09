package migrations

import (
	"context"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const schemaVersionTable = `
CREATE TABLE IF NOT EXISTS schema_versions (
    version Int64 NOT NULL,
    applied_at Timestamp NOT NULL,
    description Utf8,
    PRIMARY KEY (version)
);
`

func GetCurrentSchemaVersion(ctx context.Context, driver *ydb.Driver) (int64, error) {
	var exists bool
	err := driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		res, err := s.DescribeTable(ctx, "schema_versions")
		exists = err == nil && res.Name == "schema_versions"
		return nil
	})
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, nil
	}

	var version int64
	query := `SELECT MAX(version) FROM schema_versions;`

	err = driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, nil)
		if err != nil {
			return err
		}
		defer res.Close()

		if res.NextResultSet(ctx) && res.NextRow() {
			return res.Scan(&version)
		}
		return nil
	})

	return version, err
}

func RecordSchemaVersion(ctx context.Context, driver *ydb.Driver, version int64, description string) error {
	query := `
	DECLARE $version AS Int64;
	DECLARE $applied_at AS Timestamp;
	DECLARE $description AS Utf8;
	
	UPSERT INTO schema_versions (version, applied_at, description)
	VALUES ($version, $applied_at, $description);
	`

	return driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		_, _, err := s.Execute(ctx,
			table.DefaultTxControl(),
			query,
			table.NewQueryParameters(
				table.ValueParam("$version", types.Int64Value(version)),
				table.ValueParam("$applied_at", types.TimestampValueFromTime(time.Now())),
				table.ValueParam("$description", types.UTF8Value(description)),
			),
		)
		return err
	})
}
