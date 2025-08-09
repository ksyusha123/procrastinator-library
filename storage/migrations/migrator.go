package migrations

import (
	"context"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

type MigrationFunc func(ctx context.Context, driver *ydb.Driver) error

var migrations = map[int64]struct {
	description string
	migrate     MigrationFunc
}{
	1: {
		description: "Initial schema",
		migrate:     CreateInitialSchema,
	},
}

func Migrate(ctx context.Context, driver *ydb.Driver) error {
	currentVersion, err := GetCurrentSchemaVersion(ctx, driver)
	if err != nil {
		return fmt.Errorf("failed to get current schema version: %w", err)
	}

	if currentVersion == 0 {
		err := driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
			return s.ExecuteSchemeQuery(ctx, schemaVersionTable)
		})
		if err != nil {
			return fmt.Errorf("failed to create schema_versions table: %w", err)
		}
	}

	for version, migration := range migrations {
		if version > currentVersion {
			if err := migration.migrate(ctx, driver); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", version, err)
			}

			if err := RecordSchemaVersion(ctx, driver, version, migration.description); err != nil {
				return fmt.Errorf("failed to record schema version %d: %w", version, err)
			}
		}
	}

	return nil
}
