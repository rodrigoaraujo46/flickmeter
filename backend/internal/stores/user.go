package stores

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
)

type userStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *userStore {
	return &userStore{db}
}

func (s userStore) Read(email string, c context.Context) (user.User, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		SELECT id, email, username, avatar_url
		FROM users WHERE email = $1
		`

	var u user.User
	if err := s.db.QueryRow(ctx, query, email).Scan(
		&u.Id,
		&u.Email,
		&u.Username,
		&u.AvatarURL,
	); err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (s userStore) ReadOrCreate(u *user.User, c context.Context) (isNew bool, err error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	existing, err := s.getUserByEmailTx(ctx, tx, u.Email)
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return false, err
		}
		*u = existing
		return false, nil
	}

	if err := s.tryInsertUserTx(ctx, tx, u); err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (s userStore) getUserByEmailTx(ctx context.Context, tx pgx.Tx, email string) (user.User, error) {
	const query = `SELECT id, email, username, avatar_url FROM users WHERE email = $1`
	var u user.User

	err := tx.QueryRow(ctx, query, email).Scan(&u.Id, &u.Email, &u.Username, &u.AvatarURL)
	return u, err
}

func (s userStore) tryInsertUserTx(ctx context.Context, tx pgx.Tx, u *user.User) error {
	const query = `
        INSERT INTO users (email, username, avatar_url)
        VALUES ($1, $2, $3)
		RETURNING id;
    `

	for range 10 {
		if u.Username == "" {
			u.SetRandomUsername()
		}

		err := tx.QueryRow(ctx, query, u.Email, u.Username, u.AvatarURL).Scan(&u.Id)
		if err == nil {
			return nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key" {
			u.Username = ""
			continue
		}

		return err
	}

	return errors.New("couldn't generate a valid usernmae")
}
