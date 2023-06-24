package cmd

import (
	"context"

	"go.albinodrought.com/creamy-board/internal/db"
	"go.albinodrought.com/creamy-board/internal/log"
)

func migrate(ctx context.Context) error {
	conn, err := bootDB5(ctx)
	if err != nil {
		log.Warnf("failed to connect to DB: %v", err)
		return err
	}
	defer conn.Close(ctx)
	err = db.Migrate(ctx, conn)
	if err != nil {
		log.Warnf("failed to migrate DB: %v", err)
		return err
	}

	storage, err := bootStorage(ctx)
	if err != nil {
		log.Warnf("failed to connect to storage: %v", err)
		return err
	}
	err = storage.Boot()
	if err != nil {
		log.Warnf("failed to boot storage: %v", err)
		return err
	}

	log.Infof("migrated!")
	return nil
}
