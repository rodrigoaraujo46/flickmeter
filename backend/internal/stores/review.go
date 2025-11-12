package stores

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/movie"
)

type reviewStore struct {
	db *pgxpool.Pool
}

func NewReviewStore(db *pgxpool.Pool) *reviewStore {
	return &reviewStore{db}
}

func (r reviewStore) Create(c context.Context, review *movie.Review) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const queryReviews = `
		INSERT INTO reviews (movie_id, user_id, rating, title, review)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
	`

	if err := tx.QueryRow(ctx,
		queryReviews,
		review.MovieId,
		review.UserId,
		review.Rating,
		review.Title,
		review.Review,
	).Scan(&review.Id); err != nil {
		return err
	}

	const queryMovies = `
		INSERT INTO movies (id, total_rating, review_count, average_rating)
		VALUES ($1, $2, 1, $3)
		ON CONFLICT (id) DO UPDATE
		SET total_rating = movies.total_rating + EXCLUDED.total_rating,
    		review_count = movies.review_count + 1,
    		average_rating = (movies.total_rating + EXCLUDED.total_rating) / (movies.review_count + 1);
	`

	if _, err := tx.Exec(ctx, queryMovies, review.MovieId, review.Rating, float64(review.Rating)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r reviewStore) ReadReviews(c context.Context, movieId, page uint) (movie.Reviews, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		SELECT r.id, r.movie_id, r.title, r.rating, r.review, r.updated_at,
			   u.id AS user_id, u.username, u.avatar_url
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.movie_id = $1
		ORDER BY r.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	const limit = 10
	rows, err := r.db.Query(ctx, query, movieId, limit, limit*(page-1))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviews movie.Reviews
	for rows.Next() {
		var review movie.Review
		if err := rows.Scan(
			&review.Id,
			&review.MovieId,
			&review.Title,
			&review.Rating,
			&review.Review,
			&review.UpdatedAt,
			&review.User.Id,
			&review.User.Username,
			&review.User.AvatarURL,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r reviewStore) ReadUserReview(c context.Context, movieId, userId uint) (movie.Review, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		SELECT r.id, r.movie_id, r.title, r.rating, r.review, r.updated_at,
			   u.id AS user_id, u.username, u.avatar_url
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.movie_id = $1 AND r.user_id = $2
	`

	var review movie.Review
	if err := r.db.QueryRow(ctx, query, movieId, userId).Scan(
		&review.Id,
		&review.MovieId,
		&review.Title,
		&review.Rating,
		&review.Review,
		&review.UpdatedAt,
		&review.User.Id,
		&review.User.Username,
		&review.User.AvatarURL,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return movie.Review{}, NewErrNotFound(err)
		}
		return movie.Review{}, err
	}

	return review, nil
}

func (r reviewStore) ReadReview(c context.Context, movieId, id uint) (movie.Review, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	const query = `
		SELECT r.id, r.title, r.rating, r.review, r.updated_at,
			   u.id AS user_id, u.username, u.avatar_url
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.movie_id = $1 AND r.id = $2
	`

	var review movie.Review
	if err := r.db.QueryRow(ctx, query, movieId, id).Scan(
		&review.Id,
		&review.Title,
		&review.Rating,
		&review.Review,
		&review.UpdatedAt,
		&review.User.Id,
		&review.User.Username,
		&review.User.AvatarURL,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return movie.Review{}, NewErrNotFound(err)
		}
		return movie.Review{}, err
	}

	return review, nil
}

func (r reviewStore) Update(c context.Context, review *movie.Review) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const queryOld = `
		SELECT rating
		FROM reviews
		WHERE id = $1 AND user_id = $2 AND movie_id = $3;
	`

	var oldRating int
	if err = tx.QueryRow(ctx,
		queryOld,
		review.Id,
		review.UserId,
		review.MovieId,
	).Scan(&oldRating); err != nil {
		return err
	}

	const queryReviews = `
		UPDATE reviews
		SET title = $1, rating = $2, review = $3
		WHERE id = $4 AND user_id = $5 AND movie_id = $6;
	`

	if res, err := tx.Exec(
		ctx, queryReviews,
		review.Title, review.Rating, review.Review,
		review.Id, review.UserId, review.MovieId,
	); err != nil {
		return err
	} else if res.RowsAffected() == 0 {
		return errors.New("no review updated — not found or unauthorized")
	}

	const queryMovies = `
		UPDATE movies
        SET total_rating = total_rating - $1 + $2,
            average_rating = (total_rating - $1 + $2) / (review_count)
        WHERE id = $3
	`

	if _, err := tx.Exec(ctx, queryMovies, oldRating, review.Rating, review.MovieId); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r reviewStore) Delete(c context.Context, review *movie.Review) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const queryOld = `
		SELECT rating
		FROM reviews
		WHERE id = $1 AND user_id = $2 AND movie_id = $3;
	`

	var oldRating int
	if err = tx.QueryRow(ctx,
		queryOld,
		review.Id,
		review.UserId,
		review.MovieId,
	).Scan(&oldRating); err != nil {
		return err
	}

	const queryReviews = `
		Delete FROM reviews
		WHERE id = $1 AND user_id = $2 AND movie_id = $3;
	`

	if res, err := tx.Exec(ctx, queryReviews, review.Id, review.UserId, review.MovieId); err != nil {
		return err
	} else if res.RowsAffected() == 0 {
		return errors.New("no review deleted — not found or unauthorized")
	}

	const queryMovies = `
		UPDATE movies
		SET
			total_rating = total_rating - $1,
			review_count = review_count - 1,
			average_rating = CASE
				WHEN review_count - 1 > 0
					THEN (total_rating - $1) / (review_count - 1)
				ELSE 0
			END
		WHERE id = $2;
	`

	if _, err := tx.Exec(ctx, queryMovies, oldRating, review.MovieId); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
