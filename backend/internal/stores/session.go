package stores

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/session"
)

type sessionStore struct {
	client redis.Client
}

func NewSessionStore(client redis.Client) *sessionStore {
	return &sessionStore{client}
}

func (s sessionStore) Create(ses session.Session, c context.Context) error {
	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	json, err := json.Marshal(ses)
	if err != nil {
		return err
	}

	s.client.Set(ctx, ses.UUID, json, time.Hour)
	return nil
}

var ErrNotFound = errors.New("row not found")

func (s sessionStore) Read(key string, c context.Context) (session.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return session.Session{}, ErrNotFound
		}
		return session.Session{}, err
	}

	var ses session.Session
	if err := json.Unmarshal([]byte(res), &ses); err != nil {
		return session.Session{}, nil
	}

	return ses, nil
}

func (s sessionStore) Delete(key string, c context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
