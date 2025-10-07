package stores

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/assert"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
)

func NewPostgresClient(c config.PostgresConfig) *pgxpool.Pool {
	client, err := pgxpool.New(context.Background(), c.Address)
	assert.NoError(err, "couldn't connect to pgsql db")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
            email 	   TEXT PRIMARY KEY,
            username   TEXT NOT NULL UNIQUE,
            avatar_url TEXT
			CHECK (email <> '')
			CHECK (username <> '')
        );`,

		`CREATE TABLE IF NOT EXISTS refresh (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL
				REFERENCES users(email)
				ON DELETE CASCADE,
			CHECK (id <> ''),
			CHECK (user_id <> '')
		);`,
	}

	for _, t := range tables {
		_, err := client.Exec(ctx, t)
		assert.NoError(err, "couldnt query tables")
	}

	return client
}
