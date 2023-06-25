package cmd

import (
	"context"
	"path"

	"go.albinodrought.com/creamy-board/internal/db"
	"go.albinodrought.com/creamy-board/internal/db/migrations"
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
	err = storage.Boot(ctx)
	if err != nil {
		log.Warnf("failed to boot storage: %v", err)
		return err
	}

	migrationFiles, err := migrations.FS.ReadDir("fs")
	if err != nil {
		log.Warnf("failed to list storage migration files: %v", err)
		return err
	}

	for _, file := range migrationFiles {
		filepath := path.Join("fs", file.Name())
		handle, err := migrations.FS.Open(filepath)
		if err != nil {
			log.Warnf("failed to open storage migration file %v: %v", filepath, err)
			return err
		}
		err = storage.Write(ctx, file.Name(), handle)
		if err != nil {
			log.Warnf("failed to write storage migration file %v: %v", filepath, err)
			return err
		}
	}

	log.Infof("migrated!")
	return nil
}
