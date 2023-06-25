package cfg

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.albinodrought.com/creamy-board/internal/db/queries"
)

var (
	DB      *pgxpool.Pool
	Querier queries.Querier
)
