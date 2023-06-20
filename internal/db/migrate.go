package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"go.albinodrought.com/creamy-board/internal/db/migrations"
)

const migrationTable = "migrations"

func Migrate(ctx context.Context, conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(ctx, conn, migrationTable)
	if err != nil {
		return err
	}
	err = migrator.LoadMigrations(migrations.FS)
	if err != nil {
		return err
	}
	err = migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	return nil
}
