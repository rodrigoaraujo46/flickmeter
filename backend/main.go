package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/handlers"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/moviedb"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/stores"
)

func main() {
	c := config.MustLoadConfig()

	e := echo.New()
	e.Debug = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://web:5173", "http://web:4173"},
		AllowCredentials: true,
	}))

	setUpUserHandler(c, e)
	setUpMovieHandler(c, e)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", c.Host, c.Port)))
}

func setUpUserHandler(c config.Config, e *echo.Echo) {
	redis := stores.NewRedisClient(c.Redis)
	psql := stores.NewPostgresClient(c.Postgres)
	userHandler := handlers.NewUserHandler(
		*stores.NewSessionStore(*redis),
		*stores.NewRefreshStore(psql),
		*stores.NewUserStore(psql),
		c.Gothic,
	)

	userHandler.RegisterRoutes(e.Group("/users"))
}

func setUpMovieHandler(c config.Config, e *echo.Echo) {
	movieHandler := handlers.NewMovieHandler(moviedb.NewClient(c.MovieDB))
	movieHandler.RegisterRoutes(e.Group("/movies"))
}
