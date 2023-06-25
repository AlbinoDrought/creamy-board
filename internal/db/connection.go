package db

import (
	"context"
	"time"

	pgxpool4 "github.com/jackc/pgx/v4/pgxpool"
	pgx5 "github.com/jackc/pgx/v5"
)

// pg-gen only supports pgx v4 at the moment...
func ConnectPool4(ctx context.Context, dsn string) (conn *pgxpool4.Pool, err error) {
	const attempts = 5
	attempt := 0

	for {
		conn, err = pgxpool4.Connect(ctx, dsn)
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

// ...while tern supports pgx v5
func Connect5(ctx context.Context, dsn string) (conn *pgx5.Conn, err error) {
	const attempts = 5
	attempt := 0

	for {
		conn, err = pgx5.Connect(ctx, dsn)
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
