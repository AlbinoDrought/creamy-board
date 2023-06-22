package cmd

import (
	"context"
	"os"

	"go.albinodrought.com/creamy-board/internal/db"
	"go.albinodrought.com/creamy-board/internal/log"
)

func migrate(ctx context.Context) error {
	conn, err := db.Connect5(ctx, os.Getenv("CREAMY_BOARD_DSN"))
	if err != nil {
		log.Warnf("failed to connect: %v", err)
		return err
	}
	defer conn.Close(ctx)
	err = db.Migrate(ctx, conn)
	if err != nil {
		log.Warnf("failed to migrate: %v", err)
		return err
	}

	log.Infof("migrated!")
	return nil
}
