package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/movie"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/stores"
)

type (
	MovieClient interface {
		GetTrending(ctx context.Context, weekly bool) (movie.Movies, error)
		GetMovie(ctx context.Context, id int32) (movie.Movie, error)
		GetVideos(ctx context.Context, id int32) (movie.Videos, error)
		Search(ctx context.Context, query string) (movie.Movies, error)
	}

	MovieStore interface {
		ReadAverageRating(ctx context.Context, movieId int32) (float64, error)
	}

	ReviewStore interface {
		Create(ctx context.Context, review movie.Review) error
		ReadReviews(ctx context.Context, movieId, page int32) (movie.Reviews, error)
		ReadUserReview(ctx context.Context, movieId, userId int32) (movie.Review, error)
		ReadReview(ctx context.Context, id int32) (movie.Review, error)
		Update(ctx context.Context, review movie.Review) (movie.Review, error)
		Delete(ctx context.Context, id int32) error
	}

	movieHandler struct {
		client      MovieClient
		movieStore  MovieStore
		reviewStore ReviewStore
	}
)

func NewMovieHandler(movieClient MovieClient, movieStore MovieStore, reviewStore ReviewStore) *movieHandler {
	return &movieHandler{movieClient, movieStore, reviewStore}
}

func (h movieHandler) RegisterRoutes(g *echo.Group, protection echo.MiddlewareFunc) {
	g.GET("/:id", h.getMovie)
	g.GET("/:id/videos", h.getVideos)
	g.GET("/trending", h.getTrending)
	g.GET("/search", h.searchMovies)
	g.GET("/:id/reviews", h.getReviews)

	g.GET("/:id/reviews/me", h.getUserReview, protection)
	g.POST("/:id/reviews", h.postReview, protection)
	g.PATCH("/:id/reviews/:reviewid", h.patchReview, protection)
	g.DELETE("/:id/reviews/:reviewid", h.deleteReview, protection)
}

func (h movieHandler) getTrending(c echo.Context) error {
	weekly, _ := strconv.ParseBool(c.QueryParam("weekly"))

	movies, err := h.client.GetTrending(c.Request().Context(), weekly)
	if err != nil {
		return err
	}

	for i, movie := range movies {
		movies[i].VoteAverage, _ = h.movieStore.ReadAverageRating(c.Request().Context(), movie.Id)
	}

	return c.JSON(http.StatusOK, movies)
}

func (h movieHandler) getMovie(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}

	movie, err := h.client.GetMovie(c.Request().Context(), int32(id))
	if err != nil {
		return err
	}

	movie.VoteAverage, err = h.movieStore.ReadAverageRating(c.Request().Context(), int32(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h movieHandler) getVideos(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movid id").SetInternal(err)
	}

	videos, err := h.client.GetVideos(c.Request().Context(), int32(id))
	if err != nil {
		return err
	}

	videos.FilterTrailersAndTeasersOnYT()
	videos.SortByRelevance()

	return c.JSON(http.StatusOK, videos)
}

func (h movieHandler) getReviews(c echo.Context) error {
	movieId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid movie id").SetInternal(err)
	}

	pageStr, page := c.QueryParam("page"), int32(1)
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 32); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid page").SetInternal(err)
		} else {
			page = int32(p)
		}
	}

	reviews, err := h.reviewStore.ReadReviews(c.Request().Context(), int32(movieId), page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, reviews)
}

func (h movieHandler) getUserReview(c echo.Context) error {
	movieId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}

	review, err := h.reviewStore.ReadUserReview(c.Request().Context(), int32(movieId), MustGetUser(c).Id)
	if err != nil {
		if errors.Is(err, stores.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "review not found").SetInternal(err)
		}
		return err
	}

	return c.JSON(http.StatusOK, review)
}

func (h movieHandler) postReview(c echo.Context) error {
	movieId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}

	f := &struct {
		Title  string `json:"title"`
		Rating int32  `json:"rating"`
		Review string `json:"review"`
	}{}
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	review := movie.NewReview(f.Title, f.Rating, f.Review)
	review.UserId = MustGetUser(c).Id
	review.MovieId = int32(movieId)

	if err := h.reviewStore.Create(c.Request().Context(), *review); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, review)
}

func (h movieHandler) patchReview(c echo.Context) error {
	reviewId, err := strconv.ParseInt(c.Param("reviewid"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid review id").SetInternal(err)
	}

	ctx := c.Request().Context()

	review, err := h.reviewStore.ReadReview(ctx, int32(reviewId))
	if err != nil {
		if errors.Is(err, stores.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		return err
	}

	if review.UserId != MustGetUser(c).Id {
		return echo.ErrForbidden.SetInternal(
			errors.New("patchReview: User doesn't own this review"))
	}

	f := &struct {
		Title  string `json:"title"`
		Rating int32  `json:"rating"`
		Review string `json:"review"`
	}{}
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	review.Title = f.Title
	review.Rating = f.Rating
	review.Review = f.Review

	result, err := h.reviewStore.Update(ctx, review)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h movieHandler) deleteReview(c echo.Context) error {
	reviewId, err := strconv.ParseInt(c.Param("reviewid"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid review id").SetInternal(err)
	}

	ctx := c.Request().Context()

	review, err := h.reviewStore.ReadReview(ctx, int32(reviewId))
	if err != nil {
		if errors.Is(err, stores.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		return err
	}

	if review.UserId != MustGetUser(c).Id {
		return echo.ErrForbidden.SetInternal(
			errors.New("deleteReview: User doesn't own this review"))
	}

	if err := h.reviewStore.Delete(ctx, review.Id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h movieHandler) searchMovies(c echo.Context) error {
	movies, err := h.client.Search(c.Request().Context(), c.QueryParam("query"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}
