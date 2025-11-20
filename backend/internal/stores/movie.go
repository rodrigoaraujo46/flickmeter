package stores

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/db"
)

type movieStore struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewMovieStore(db *pgxpool.Pool) *movieStore {
	return &movieStore{db, time.Second}
}

func (s movieStore) ReadAverageRating(c context.Context, id int32) (float64, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)

	movie, err := q.ReadMovie(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return float64(movie.TotalRating) / float64(movie.ReviewCount), nil
}
