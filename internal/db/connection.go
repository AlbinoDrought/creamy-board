package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, dsn string) (conn *pgx.Conn, err error) {
	const attempts = 5
	attempt := 0

	for {
		conn, err = pgx.Connect(ctx, dsn)
		if err == nil {
			return
		}
		attempt++
		if attempt > attempts {
			return
		}
		time.Sleep(time.Second)
	}
}
