package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/movie"
)

type MovieClient interface {
	GetTrending(ctx context.Context, weekly bool) (movie.Movies, error)
}

type movieHandler struct {
	client MovieClient
}

func NewMovieHandler(movieClient MovieClient) *movieHandler {
	return &movieHandler{movieClient}
}

func (h movieHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/trending", h.getTrending)
}

func (h movieHandler) getTrending(c echo.Context) error {
	weekly, _ := strconv.ParseBool(c.QueryParam("weekly"))

	movies, err := h.client.GetTrending(c.Request().Context(), weekly)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}
