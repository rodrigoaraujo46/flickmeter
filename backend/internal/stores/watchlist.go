package stores

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type watchlistStore struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewWatchlistStore(db *pgxpool.Pool) *watchlistStore {
	return &watchlistStore{db, time.Second}
}
