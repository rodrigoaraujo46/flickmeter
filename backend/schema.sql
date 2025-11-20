CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE CHECK (length(username) >= 5 AND length(username) <= 30),
    email TEXT NOT NULL UNIQUE CHECK (length(email) >= 3 AND length(email) <= 254),
    avatar_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS movies (
    id INT PRIMARY KEY,
    total_rating INT NOT NULL DEFAULT 0,
    review_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (movie_id, user_id),
    title TEXT NOT NULL CHECK (length(title) >= 1 AND length(title) <= 100),
    rating INT NOT NULL CHECK (rating >= 0 AND rating <= 10),
    review TEXT NOT NULL CHECK (length(review) >= 1 AND length(review) <= 1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE watchlists (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id INT NOT NULL,
    watched BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, movie_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_movie_updated_at
ON reviews (movie_id, updated_at DESC);


DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $func$
        BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
        END;
        $func$ LANGUAGE plpgsql';
    END IF;
END
$$;

DO $$
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

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'watchlists_updated_at') THEN
        EXECUTE 'CREATE TRIGGER watchlists_updated_at
        BEFORE UPDATE ON watchlists
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
    END IF;
END
$$;
