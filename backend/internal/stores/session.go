package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	return s.client.Set(ctx, ses.UUID, json, time.Hour).Err()
}

var ErrNotFound = errors.New("row not found")

func NewErrNotFound(err error) error {
	return fmt.Errorf("%w: %v", ErrNotFound, err)
}

func (s sessionStore) ReadAndRefresh(key string, c context.Context) (session.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := s.client.GetEx(ctx, key, time.Hour).Result()
	if err != nil {
		if err == redis.Nil {
			return session.Session{}, NewErrNotFound(err)
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

	return s.client.Del(ctx, key).Err()
}
