package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/handlers"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/movieapi"
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

	setUpHandlers(c, e)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", c.Host, c.Port)))
}

func setUpHandlers(c config.Config, e *echo.Echo) {
	psql := stores.NewPostgresClient(c.Postgres)
	redis := stores.NewRedisClient(c.Redis)

	userHandler := handlers.NewUserHandler(
		stores.NewSessionStore(*redis),
		stores.NewRefreshStore(psql),
		stores.NewUserStore(psql),
		c.Gothic,
	)
	movieHandler := handlers.NewMovieHandler(
		movieapi.NewClient(c.MovieAPI),
		stores.NewMovieStore(psql),
		stores.NewReviewStore(psql),
	)

	e.Use(userHandler.AuthMiddleware)
	userHandler.RegisterRoutes(e.Group("/users"))
	movieHandler.RegisterRoutes(e.Group("/movies"))
}
