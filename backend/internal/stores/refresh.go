package stores

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/refresh"
)

type refreshStore struct {
	db *pgxpool.Pool
}

func NewRefreshStore(db *pgxpool.Pool) *refreshStore {
	return &refreshStore{db}
}

func (s refreshStore) Create(refresh refresh.Refresh, c context.Context) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
        INSERT INTO refresh (id, user_id)
        VALUES ($1, $2)
	`

	_, err := s.db.Exec(ctx, query, refresh.UUID, refresh.User.Id)

	return err
}

func (s refreshStore) Read(uuid string, c context.Context) (refresh.Refresh, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		SELECT u.id, u.email, u.username, u.avatar_url
        FROM refresh r
        JOIN users u ON r.user_id = u.id
        WHERE r.id = $1
	`

	var ref refresh.Refresh
	ref.UUID = uuid
	if err := s.db.QueryRow(ctx, query, uuid).Scan(&ref.User.Id, &ref.User.Email, &ref.User.Username, &ref.User.AvatarURL); err != nil {
		if err == pgx.ErrNoRows {
			return refresh.Refresh{}, NewErrNotFound(err)
		}
		return refresh.Refresh{}, err
	}

	return ref, nil
}

func (s refreshStore) Delete(uuid string, c context.Context) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		DELETE FROM refresh
		WHERE id = $1;
	`

	if _, err := s.db.Exec(ctx, query, uuid); err != nil {
		return err
	}

	return nil
}
