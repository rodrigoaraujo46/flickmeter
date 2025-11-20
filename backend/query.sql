-- name: ReadMovie :one
SELECT * FROM movies
WHERE id = $1;

-- name: CreateRefresh :exec
INSERT INTO refresh (id, user_id)
VALUES ($1, $2);

-- name: ReadRefresh :one
SELECT sqlc.embed(r), sqlc.embed(u)
FROM refresh r
JOIN users u ON r.user_id = u.id
WHERE r.id = $1;

-- name: DeleteRefresh :exec
DELETE FROM refresh
WHERE id = $1;

-- name: CreateReview :exec
INSERT INTO reviews (movie_id, user_id, rating, title, review)
VALUES ($1, $2, $3, $4, $5);

-- name: IncrementMovieRating :exec
INSERT INTO movies (id, total_rating, review_count)
VALUES ($1, sqlc.arg(rating), 1)
    ON CONFLICT (id) DO UPDATE
    SET total_rating = movies.total_rating + sqlc.arg(rating),
    review_count = movies.review_count + 1;

-- name: DecrementMovieRating :exec
UPDATE movies
    SET total_rating = movies.total_rating - sqlc.arg(rating),
    review_count = movies.review_count - 1
WHERE id = $1;

-- name: UpdateMovieRating :exec
UPDATE movies
SET total_rating = total_rating - sqlc.arg(old_rating) + sqlc.arg(new_rating)
WHERE id = $1;

-- name: ReadReviews :many
SELECT sqlc.embed(reviews), sqlc.embed(users)
FROM reviews
JOIN users ON reviews.user_id = users.id
WHERE reviews.movie_id = $1
ORDER BY reviews.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ReadUserReview :one
SELECT sqlc.embed(reviews), sqlc.embed(users)
FROM reviews
JOIN users ON reviews.user_id = users.id
WHERE reviews.movie_id = $1 AND reviews.user_id = $2;

-- name: ReadReview :one
SELECT sqlc.embed(reviews), sqlc.embed(users)
FROM reviews
JOIN users ON reviews.user_id = users.id
WHERE reviews.id = $1;

-- name: UpdateReview :one
UPDATE reviews
SET title = $2, rating = $3, review = $4
FROM users
WHERE reviews.id = $1
RETURNING sqlc.embed(reviews), sqlc.embed(users);

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;

-- name: ReadUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: ReadUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, avatar_url)
VALUES ($1, $2, $3)
RETURNING *;
