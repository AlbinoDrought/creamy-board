package cmd

import (
	"context"
	"net/http"
	"os"

	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/db/queries"
	"go.albinodrought.com/creamy-board/internal/log"
	"go.albinodrought.com/creamy-board/internal/web"
)

func serve(ctx context.Context) error {
	var err error
	cfg.DB, err = bootDB4(ctx)
	if err != nil {
		log.Warnf("failed to connect to db: %v", err)
		return err
	}
	cfg.Querier = queries.NewQuerier(cfg.DB)

	cfg.Storage, err = bootStorage(ctx)
	if err != nil {
		log.Warnf("failed to connect to storage: %v", err)
		return err
	}

	listenAddr := os.Getenv("CREAMY_LISTEN_ADDRESS")
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	err = http.ListenAndServe(listenAddr, web.Router())
	if err != nil {
		log.Warnf("error during listen: %v", err)
		return err
	}

	return nil
}
