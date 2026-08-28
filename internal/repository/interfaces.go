package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
