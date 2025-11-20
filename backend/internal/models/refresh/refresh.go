package refresh

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
)

type Refresh struct {
	UUID    uuid.UUID
	User    user.User
	Expires time.Time
}

func New(uuid uuid.UUID, user user.User, keep bool) *Refresh {
	var expires time.Time
	if keep {
		expires = time.Now().Add(720 * time.Hour)
	}

	return &Refresh{UUID: uuid, User: user, Expires: expires}
}

func (r Refresh) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:     "refresh",
		Path:     "/",
		Value:    r.UUID.String(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  r.Expires,
	}
}
