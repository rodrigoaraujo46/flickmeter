package movieapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

func NewClient(c config.MovieAPI) *client {
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

	if res.StatusCode != http.StatusOK {
		var jsonMap map[string]any
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return nil, err
		}
		return nil, echo.NewHTTPError(res.StatusCode, jsonMap["status_message"])
	}

	var movieRes struct {
		Results movie.Movies `json:"results"`
	}

	err = json.Unmarshal(body, &movieRes)
	if err != nil {
		return nil, err
	}

	return movieRes.Results, nil
}

func (c client) GetMovie(ctx context.Context, id uint) (movie.Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?language=en-US", id)

	var empty movie.Movie
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return empty, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return empty, err
	}

	defer func() {
		err := res.Body.Close()
		assert.NoError(err, "movieAPI.GetMovie")
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	if res.StatusCode != http.StatusOK {
		var jsonMap map[string]any
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return empty, err
		}
		return empty, echo.NewHTTPError(res.StatusCode, jsonMap["status_message"])
	}

	var movie movie.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return empty, err
	}

	return movie, nil
}

func (c client) GetVideos(ctx context.Context, id uint) (movie.Videos, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/videos?language=en-US", id)

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
		assert.NoError(err, "movieAPI.GetMovie")
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var jsonMap map[string]any
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return nil, err
		}
		return nil, echo.NewHTTPError(res.StatusCode, jsonMap["status_message"])
	}

	var videoRes struct {
		Results movie.Videos `json:"results"`
	}

	err = json.Unmarshal(body, &videoRes)
	if err != nil {
		return nil, err
	}

	return videoRes.Results, nil
}

func (c client) Search(ctx context.Context, query string) (movie.Movies, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?query=%s", query)

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
		assert.NoError(err, "movieAPI.Search")
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var jsonMap map[string]any
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return nil, err
		}
		return nil, echo.NewHTTPError(res.StatusCode, jsonMap["status_message"])
	}

	var movieRes struct {
		Results movie.Movies `json:"results"`
	}

	err = json.Unmarshal(body, &movieRes)
	if err != nil {
		return nil, err
	}

	return movieRes.Results, nil
}
