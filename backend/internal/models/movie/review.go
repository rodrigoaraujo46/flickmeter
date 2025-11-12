package movie

import (
	"time"

	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
)

type Review struct {
	Id        uint      `json:"id"`
	MovieId   uint      `json:"movie_id"`
	UserId    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Rating    uint      `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      user.User `json:"user"`
}

type Reviews []Review
