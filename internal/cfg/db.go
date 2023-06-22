package cfg

import (
	"github.com/jackc/pgx/v4"
	"go.albinodrought.com/creamy-board/internal/db/queries"
)

var (
	DB      *pgx.Conn
	Querier queries.Querier
)
