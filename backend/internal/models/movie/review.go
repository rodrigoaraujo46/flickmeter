package movie

import (
	"time"

	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
)

type Review struct {
	Id        int32     `json:"id"`
	MovieId   int32     `json:"movie_id"`
	UserId    int32     `json:"user_id"`
	Title     string    `json:"title"`
	Rating    int32     `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      user.User `json:"user"`
}

type Reviews []Review

func NewReview(title string, rating int32, review string) *Review {
	return &Review{Title: title, Rating: rating, Review: review}
}
