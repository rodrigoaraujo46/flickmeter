package stores

import (
	"github.com/redis/go-redis/v9"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
)

func NewRedisClient(c config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: c.Address})
}
