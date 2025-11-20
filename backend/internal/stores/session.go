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

var ErrNotFound = errors.New("row not found")

func NewErrNotFound(err error) error {
	return fmt.Errorf("%w: %v", ErrNotFound, err)
}

type sessionStore struct {
	client  redis.Client
	timeout time.Duration
}

func NewSessionStore(client redis.Client) *sessionStore {
	return &sessionStore{client, time.Second}
}

func (s sessionStore) Create(c context.Context, ses session.Session) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	json, err := json.Marshal(ses)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, ses.UUID, json, time.Hour).Err()
}

func (s sessionStore) ReadAndRefresh(c context.Context, key string) (session.Session, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
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

func (s sessionStore) Delete(c context.Context, key string) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	return s.client.Del(ctx, key).Err()
}
