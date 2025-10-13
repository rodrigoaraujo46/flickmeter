package moviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rodrigoaraujo46/assert"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/movie"
)

type client struct {
	http *http.Client
}

type authTransport struct {
	base  http.RoundTripper
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.base.RoundTrip(req)
}

func NewClient(c config.MovieDB) *client {
	httpClient := &http.Client{
		Timeout:   time.Second,
		Transport: &authTransport{base: http.DefaultTransport, token: c.Token},
	}

	return &client{
		http: httpClient,
	}
}

func (c client) GetTrending(ctx context.Context, weekly bool) (movie.Movies, error) {
	url := "https://api.themoviedb.org/3/trending/movie/day?language=en-US"
	if weekly {
		url = "https://api.themoviedb.org/3/trending/movie/week?language=en-US"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := res.Body.Close()
		assert.NoError(err, "movieAPI.GetTrending")
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var movieRes struct {
		Results movie.Movies `json:"results"`
	}

	err = json.Unmarshal(body, &movieRes)

	return movieRes.Results, err
}
