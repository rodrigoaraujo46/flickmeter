package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rodrigoaraujo46/assert"
)

type Config struct {
	Host     string
	Port     string
	Redis    Redis
	Postgres Postgres
	Gothic   Gothic
	MovieDB  MovieDB
}

type Redis struct {
	Address string
}

type Postgres struct {
	Address string
}

type MovieDB struct {
	Token string
}

type Gothic struct {
	Providers      map[string]oAuthProvider
	CookieStoreKey string
}

type oAuthProvider struct {
	Client   string
	Secret   string
	Callback string
}

func MustLoadConfig() Config {
	assert.NoError(godotenv.Load(), "Couldn't open .env files")

	return Config{
		Host: mustLoadEnv("HOST"),
		Port: mustLoadEnv("PORT"),
		Redis: Redis{
			Address: mustLoadEnv("REDIS_ADDR"),
		},
		Postgres: Postgres{
			Address: mustLoadEnv("POSTGRES_ADDR"),
		},
		Gothic: Gothic{
			CookieStoreKey: mustLoadEnv("COOKIE_STORE_KEY"),
			Providers:      mustLoadProviders(),
		},
		MovieDB: MovieDB{
			Token: mustLoadEnv("MOVIE_DB_TOKEN"),
		},
	}
}

func mustLoadEnv(name string) string {
	port, found := os.LookupEnv(name)
	assert.Assert(found, fmt.Sprintf("No %s in .env", name))

	return port
}

func mustLoadProviders() map[string]oAuthProvider {
	names := strings.Split(mustLoadEnv("PROVIDERS"), ",")
	configs := make(map[string]oAuthProvider, len(names))

	for _, name := range names {
		config := oAuthProvider{
			Client:   mustLoadEnv(fmt.Sprintf("%s_CLIENT", strings.ToUpper(name))),
			Secret:   mustLoadEnv(fmt.Sprintf("%s_SECRET", strings.ToUpper(name))),
			Callback: fmt.Sprintf("http://localhost:5173/api/users/auth/%s/callback", name),
		}
		configs[name] = config
	}

	return configs
}
