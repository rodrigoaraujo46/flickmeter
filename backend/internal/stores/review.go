package stores

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/db"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/movie"
)

type reviewStore struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewReviewStore(db *pgxpool.Pool) *reviewStore {
	return &reviewStore{db, time.Second}
}

func (s reviewStore) Create(c context.Context, review movie.Review) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := db.New(s.db).WithTx(tx)

	if err := qtx.CreateReview(ctx, db.CreateReviewParams{
		MovieID: review.MovieId,
		UserID:  review.UserId,
		Rating:  review.Rating,
		Title:   review.Title,
		Review:  review.Review,
	}); err != nil {
		return err
	}

	if err := qtx.IncrementMovieRating(ctx, db.IncrementMovieRatingParams{
		ID:     review.MovieId,
		Rating: review.Rating,
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s reviewStore) ReadReviews(c context.Context, movieId, page int32) (movie.Reviews, error) {
	const limit = 10

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)
	results, err := q.ReadReviews(ctx, db.ReadReviewsParams{MovieID: movieId, Limit: limit, Offset: page})
	if err != nil {
		return nil, err
	}

	reviews := make(movie.Reviews, 0, len(results))
	for _, r := range results {
		review := reviewRowToReview(r.Review)
		review.User = userRowToUser(r.User)
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (s reviewStore) ReadUserReview(c context.Context, movieId, userId int32) (review movie.Review, err error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)

	result, err := q.ReadUserReview(ctx, db.ReadUserReviewParams{MovieID: movieId, UserID: userId})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return review, NewErrNotFound(err)
		}
		return movie.Review{}, err
	}

	review = reviewRowToReview(result.Review)
	review.User = userRowToUser(result.User)
	return review, nil
}

func (s reviewStore) ReadReview(c context.Context, id int32) (review movie.Review, err error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	q := db.New(s.db)

	result, err := q.ReadReview(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return review, NewErrNotFound(err)
		}
		return review, err
	}

	review = reviewRowToReview(result.Review)
	review.User = userRowToUser(result.User)
	return review, nil
}

func (s reviewStore) Update(c context.Context, review movie.Review) (movie.Review, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return review, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := db.New(s.db).WithTx(tx)

	oldReview, err := qtx.ReadReview(ctx, review.Id)
	if err != nil {
		return review, err
	}

	result, err := qtx.UpdateReview(ctx, db.UpdateReviewParams{
		ID: review.Id, Title: review.Title,
		Rating: review.Rating, Review: review.Review,
	})
	if err != nil {
		return review, err
	}

	if err := qtx.UpdateMovieRating(ctx, db.UpdateMovieRatingParams{
		ID:        review.MovieId,
		OldRating: oldReview.Review.Rating, NewRating: result.Review.Rating,
	}); err != nil {
		return review, err
	}

	if err := tx.Commit(ctx); err != nil {
		return review, err
	}

	review = reviewRowToReview(result.Review)
	review.User = userRowToUser(result.User)
	return review, err
}

func (s reviewStore) Delete(c context.Context, id int32) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(ctx) }()

	qtx := db.New(s.db).WithTx(tx)

	oldReview, err := qtx.ReadReview(ctx, id)
	if err != nil {
		return err
	}

	if err := qtx.DeleteReview(ctx, id); err != nil {
		return err
	}

	if err := qtx.DecrementMovieRating(ctx, db.DecrementMovieRatingParams{
		ID:     id,
		Rating: oldReview.Review.Rating,
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func reviewRowToReview(r db.Review) movie.Review {
	return movie.Review{
		Id:        r.ID,
		MovieId:   r.MovieID,
		UserId:    r.UserID,
		Title:     r.Title,
		Rating:    r.Rating,
		Review:    r.Review,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
