package stores

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/db"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/refresh"
)

type refreshStore struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewRefreshStore(db *pgxpool.Pool) *refreshStore {
	return &refreshStore{db, time.Second}
}

func (s refreshStore) Create(c context.Context, refresh refresh.Refresh) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)

	return q.CreateRefresh(ctx, db.CreateRefreshParams{
		ID:     refresh.UUID,
		UserID: refresh.User.Id,
	})
}

func (s refreshStore) Read(c context.Context, uuid uuid.UUID) (refresh refresh.Refresh, err error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)
	result, err := q.ReadRefresh(ctx, uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return refresh, NewErrNotFound(err)
		}
		return refresh, err
	}

	refresh.UUID, refresh.User = result.Refresh.ID, userRowToUser(result.User)
	return refresh, nil
}

func (s refreshStore) Delete(c context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)
	return q.DeleteRefresh(ctx, uuid)
}
