package session

import (
	"net/http"

	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
)

type Session struct {
	UUID string    `json:"uuid"`
	User user.User `json:"user"`
}

func New(uuid string, user user.User) *Session {
	return &Session{
		UUID: uuid,
		User: user,
	}
}

func (r Session) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Path:     "/",
		Value:    r.UUID,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
