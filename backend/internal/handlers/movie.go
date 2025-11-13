package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/movie"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/stores"
)

type MovieClient interface {
	GetTrending(ctx context.Context, weekly bool) (movie.Movies, error)
	GetMovie(ctx context.Context, id uint) (movie.Movie, error)
	GetVideos(ctx context.Context, id uint) (movie.Videos, error)
	Search(ctx context.Context, query string) (movie.Movies, error)
}

type MovieStore interface {
	ReadAverageRating(ctx context.Context, movieId uint) (float64, error)
}

type ReviewStore interface {
	Create(ctx context.Context, review *movie.Review) error
	ReadReviews(ctx context.Context, movieId, page uint) (movie.Reviews, error)
	ReadUserReview(ctx context.Context, movieId, userId uint) (movie.Review, error)
	Update(ctx context.Context, review *movie.Review) error
	Delete(ctx context.Context, review *movie.Review) error
}

type movieHandler struct {
	client      MovieClient
	movieStore  MovieStore
	reviewStore ReviewStore
}

func NewMovieHandler(movieClient MovieClient, movieStore MovieStore, reviewStore ReviewStore) *movieHandler {
	return &movieHandler{movieClient, movieStore, reviewStore}
}

func (h movieHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/trending", h.getTrending)

	g.GET("/:id", h.getMovie)

	g.GET("/:id/videos", h.getVideos)

	g.GET("/:id/reviews", h.getReviews)
	g.GET("/:id/reviews/me", h.getUserReview)
	g.POST("/:id/reviews", h.postReview)
	g.PATCH("/:id/reviews/:reviewid", h.patchReview)
	g.DELETE("/:id/reviews/:reviewid", h.deleteReview)
	g.GET("/search", h.searchMovies)
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}

	movie, err := h.client.GetMovie(c.Request().Context(), uint(id))
	if err != nil {
		return err
	}

	movie.VoteAverage, err = h.movieStore.ReadAverageRating(c.Request().Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h movieHandler) getVideos(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movid id").SetInternal(err)
	}

	videos, err := h.client.GetVideos(c.Request().Context(), uint(id))
	if err != nil {
		return err
	}

	videos.FilterTrailersAndTeasersOnYT()
	videos.SortByRelevance()

	return c.JSON(http.StatusOK, videos)
}

func (h movieHandler) getReviews(c echo.Context) error {
	movieId, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid movie id").SetInternal(err)
	}

	pageStr, page := c.QueryParam("page"), uint(1)
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 0); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid page").SetInternal(err)
		} else {
			page = uint(p)
		}
	}

	reviews, err := h.reviewStore.ReadReviews(c.Request().Context(), uint(movieId), page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, reviews)
}

func (h movieHandler) getUserReview(c echo.Context) error {
	movieId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}

	user, ok := c.Get("user").(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "no user found")
	}

	review, err := h.reviewStore.ReadUserReview(c.Request().Context(), uint(movieId), user.Id)
	if err != nil {
		if errors.Is(err, stores.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "review not found").SetInternal(err)
		}
		return err
	}

	return c.JSON(http.StatusOK, review)
}

func (h movieHandler) postReview(c echo.Context) error {
	movieId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}

	user, ok := c.Get("user").(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "no user found").
			SetInternal(errors.New("couldn't assert 'user' to type User"))
	}

	review := new(movie.Review)
	if err := c.Bind(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	review.UserId = user.Id
	review.MovieId = uint(movieId)

	if err := h.reviewStore.Create(c.Request().Context(), review); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, review)
}

func (h movieHandler) patchReview(c echo.Context) error {
	review := new(movie.Review)

	movieId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}
	review.MovieId = uint(movieId)

	reviewId, err := strconv.ParseUint(c.Param("reviewid"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid review id").SetInternal(err)
	}
	review.Id = uint(reviewId)

	user, ok := c.Get("user").(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "no user found").
			SetInternal(errors.New("couldn't assert 'user' to type User"))
	}
	review.UserId = user.Id

	if err := c.Bind(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	if err := h.reviewStore.Update(c.Request().Context(), review); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h movieHandler) deleteReview(c echo.Context) error {
	review := new(movie.Review)

	movieId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid movie id").SetInternal(err)
	}
	review.MovieId = uint(movieId)

	reviewId, err := strconv.ParseUint(c.Param("reviewid"), 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not a valid review id").SetInternal(err)
	}
	review.Id = uint(reviewId)

	user, ok := c.Get("user").(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "no user found").
			SetInternal(errors.New("couldn't assert 'user' to type User"))
	}
	review.UserId = user.Id

	if err := c.Bind(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	if err := h.reviewStore.Delete(c.Request().Context(), review); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h movieHandler) searchMovies(c echo.Context) error {
	movies, err := h.client.Search(c.Request().Context(), c.QueryParam("query"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}
