package stores

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/assert"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
)

func NewPostgresClient(c config.Postgres) *pgxpool.Pool {
	client, err := pgxpool.New(context.Background(), c.Address)
	assert.NoError(err, "couldn't connect to pgsql db")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE CHECK (length(username) >= 5 AND length(username) <= 30),
			email TEXT NOT NULL UNIQUE CHECK (length(email) >= 3 AND length(email) <= 254),
			avatar_url TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,

		`CREATE TABLE IF NOT EXISTS refresh (
			id UUID PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,

		`CREATE TABLE IF NOT EXISTS movies (
			id INT PRIMARY KEY,
			total_rating INT NOT NULL DEFAULT 0,
		    review_count INT NOT NULL DEFAULT 0,
			average_rating FLOAT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,

		`CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			movie_id INT NOT NULL,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE (movie_id, user_id),
			title TEXT NOT NULL CHECK (length(title) >= 1 AND length(title) <= 100),
			rating INT CHECK (rating >= 0 AND rating <= 10),
			review TEXT NOT NULL CHECK (length(review) >= 1 AND length(review) <= 1000),
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,

		`CREATE INDEX IF NOT EXISTS idx_reviews_movie_updated_at
		 	ON reviews (movie_id, updated_at DESC);
		`,

		`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
					EXECUTE 'CREATE OR REPLACE FUNCTION update_updated_at_column()
							 RETURNS TRIGGER AS $func$
							 BEGIN
								 NEW.updated_at = now();
								 RETURN NEW;
							 END;
							 $func$ LANGUAGE plpgsql';
				END IF;
			END
		$$;`,

		`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'users_updated_at') THEN
					EXECUTE 'CREATE TRIGGER users_updated_at
							 BEFORE UPDATE ON users
							 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
				END IF;

				IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'refresh_updated_at') THEN
					EXECUTE 'CREATE TRIGGER refresh_updated_at
							 BEFORE UPDATE ON refresh
							 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
				END IF;

				IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'movies_updated_at') THEN
					EXECUTE 'CREATE TRIGGER movies_updated_at
							 BEFORE UPDATE ON movies
							 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
				END IF;

				IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'reviews_updated_at') THEN
					EXECUTE 'CREATE TRIGGER reviews_updated_at
							 BEFORE UPDATE ON reviews
							 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
				END IF;
			END
		$$;`,
	}

	for _, t := range queries {
		_, err := client.Exec(ctx, t)
		assert.NoError(err, "couldnt query tables")
	}

	return client
}
