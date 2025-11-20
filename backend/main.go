package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rodrigoaraujo46/assert"
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
		AllowOrigins: []string{"http://web:5173", "http://web:4173"}, AllowCredentials: true,
	}))

	setUpHandlers(c, e)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", c.Host, c.Port)))
}

func setUpHandlers(c config.Config, e *echo.Echo) {
	psql, err := pgxpool.New(context.Background(), c.Postgres.Address)
	assert.NoError(err, "failed to connect to postgres")

	ExecSchema(psql)

	redis := redis.NewClient(&redis.Options{Addr: c.Redis.Address})

	userHandler := handlers.NewUserHandler(stores.NewSessionStore(*redis),
		stores.NewRefreshStore(psql), stores.NewUserStore(psql), c.Gothic)

	movieHandler := handlers.NewMovieHandler(movieapi.NewClient(c.MovieAPI),
		stores.NewMovieStore(psql), stores.NewReviewStore(psql))

	watchlistHandler := handlers.NewWatchlistHandler(stores.NewWatchlistStore(psql))

	userHandler.RegisterRoutes(e.Group("/users"), userHandler.Protection)
	movieHandler.RegisterRoutes(e.Group("/movies"), userHandler.Protection)
	watchlistHandler.RegisterRoutes(e.Group("/watchlists"))
}

func ExecSchema(db *pgxpool.Pool) error {
	ctx := context.Background()

	sqlBytes, err := os.ReadFile("./schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	sql := string(sqlBytes)

	_, err = db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to exec schema: %w", err)
	}

	return nil
}
