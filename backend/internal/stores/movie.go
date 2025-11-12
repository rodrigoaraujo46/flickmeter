package stores

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type movieStore struct {
	db *pgxpool.Pool
}

func NewMovieStore(db *pgxpool.Pool) *movieStore {
	return &movieStore{db}
}

func (m movieStore) ReadAverageRating(c context.Context, id uint) (float64, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		SELECT average_rating
		FROM movies
		WHERE id = $1;
	`

	var average float64
	if err := m.db.QueryRow(ctx, query, id).Scan(&average); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return average, nil
}
