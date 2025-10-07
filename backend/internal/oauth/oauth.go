package oauth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	githubProvider "github.com/markbates/goth/providers/github"
	googleProvider "github.com/markbates/goth/providers/google"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
)

func StartOAuth(conf config.GothicConfig) {
	gothic.Store = sessions.NewCookieStore([]byte(conf.CookieStoreKey))

	google := conf.Providers["google"]
	github := conf.Providers["github"]
	goth.UseProviders(
		googleProvider.New(google.Client, google.Secret, google.Callback, "profile", "email"),
		githubProvider.New(github.Client, github.Secret, github.Callback, "user:email"),
	)
}
