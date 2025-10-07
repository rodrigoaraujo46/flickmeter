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
	Redis    RedisConfig
	Postgres PostgresConfig
	Gothic   GothicConfig
}

type RedisConfig struct {
	Address string
}

type PostgresConfig struct {
	Address string
}

type GothicConfig struct {
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
		Redis: RedisConfig{
			Address: mustLoadEnv("REDIS_ADDR"),
		},
		Postgres: PostgresConfig{
			Address: mustLoadEnv("POSTGRES_ADDR"),
		},
		Gothic: GothicConfig{
			CookieStoreKey: mustLoadEnv("COOKIE_STORE_KEY"),
			Providers:      mustLoadProviders(),
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
			Callback: fmt.Sprintf("http://localhost:5173/api/user/auth/%s/callback", name),
		}
		configs[name] = config
	}

	return configs
}
