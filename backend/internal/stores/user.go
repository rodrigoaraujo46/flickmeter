package stores

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/assert"
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
		SELECT email, username, avatar_url
		FROM users WHERE email = $1
		`

	var u user.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&u.Email,
		&u.Username,
		&u.AvatarURL,
	)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (s userStore) ReadOrCreate(u *user.User, c context.Context) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if tx.Conn().IsClosed() {
			err := tx.Rollback(ctx)
			assert.NoError(err, "Transaction already closed")
		}
	}()

	existing, err := s.getUserByEmailTx(ctx, tx, u.Email)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("gg")
		return err
	}
	if err == nil {
		*u = existing
		return tx.Commit(ctx)
	}

	if err := s.tryInsertUserTx(ctx, tx, u); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s userStore) getUserByEmailTx(ctx context.Context, tx pgx.Tx, email string) (user.User, error) {
	const query = `SELECT email, username, avatar_url FROM users WHERE email = $1`
	var u user.User

	err := tx.QueryRow(ctx, query, email).Scan(&u.Email, &u.Username, &u.AvatarURL)
	return u, err
}

func (s userStore) tryInsertUserTx(ctx context.Context, tx pgx.Tx, u *user.User) error {
	const query = `
        INSERT INTO users (email, username, avatar_url)
        VALUES ($1, $2, $3)
    `

	for range 10 {
		if u.Username == "" {
			u.SetRandomUsername()
		}

		_, err := tx.Exec(ctx, query, u.Email, u.Username, u.AvatarURL)
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
