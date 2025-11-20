package oauth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
)

func StartOAuth(conf config.Gothic) {
	gothic.Store = sessions.NewCookieStore([]byte(conf.CookieStoreKey))
	googleProvider, githubProvider := conf.Providers["google"], conf.Providers["github"]
	goth.UseProviders(
		google.New(googleProvider.Client, googleProvider.Secret, googleProvider.Callback, "profile", "email"),
		github.New(githubProvider.Client, githubProvider.Secret, githubProvider.Callback, "user:email"),
	)
}
