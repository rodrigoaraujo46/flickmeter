package handlers

import "github.com/labstack/echo/v4"

type (
	WatchlistStore interface{}

	watchlistHandler struct {
		watchlistStore WatchlistStore
	}
)

func NewWatchlistHandler(watchlistStore WatchlistStore) *watchlistHandler {
	return &watchlistHandler{watchlistStore}
}

func (w *watchlistHandler) RegisterRoutes(group *echo.Group) {
}
