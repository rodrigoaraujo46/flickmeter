package stores

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/db"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
)

type userStore struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewUserStore(db *pgxpool.Pool) *userStore {
	return &userStore{db, time.Second}
}

func (s userStore) Read(c context.Context, id int32) (user.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)

	u, err := q.ReadUser(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	return userRowToUser(u), nil
}

func (s userStore) ReadOrCreate(c context.Context, u user.User) (user user.User, isNew bool, err error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return user, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := db.New(s.db).WithTx(tx)

	existing, err := qtx.ReadUserByEmail(ctx, u.Email)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return user, false, err
		}
		return userRowToUser(existing), false, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return user, false, err
	}

	created, err := s.createUserWithRetry(ctx, qtx, u)
	if err != nil {
		return user, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return user, false, err
	}
	return userRowToUser(created), true, nil
}

func (s userStore) createUserWithRetry(ctx context.Context, qtx *db.Queries, u user.User) (db.User, error) {
	for range 10 {
		if u.Username == "" {
			u.SetRandomUsername()
		}

		res, err := qtx.CreateUser(ctx, db.CreateUserParams{
			Username:  u.Username,
			Email:     u.Email,
			AvatarUrl: u.AvatarURL,
		})
		if err == nil {
			return res, nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "users_username_key" {
			u.Username = ""
			continue
		}

		return db.User{}, err
	}

	return db.User{}, errors.New("failed to generate a valid username")
}

func userRowToUser(u db.User) user.User {
	return user.User{
		Id:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		AvatarURL: u.AvatarUrl,
	}
}
