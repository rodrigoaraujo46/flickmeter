package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/config"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/refresh"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/session"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/models/user"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/oauth"
	"github.com/rodrigoaraujo46/flickmeter/backend/internal/stores"
)

type SessionStore interface {
	Create(session session.Session, ctx context.Context) error
	ReadAndRefresh(uuid string, ctx context.Context) (session session.Session, err error)
	Delete(uuid string, ctx context.Context) error
}

type RefreshStore interface {
	Create(refreshToken refresh.Refresh, ctx context.Context) error
	Read(uuid string, ctx context.Context) (refresh.Refresh, error)
	Delete(uuid string, ctx context.Context) error
}

type UserStore interface {
	ReadOrCreate(user *user.User, ctx context.Context) (isNew bool, err error)
	Read(id string, ctx context.Context) (user user.User, err error)
}

type userHandler struct {
	sessionStore SessionStore
	refreshStore RefreshStore
	userStore    UserStore
}

func NewUserHandler(authStore SessionStore, refreshStore RefreshStore, userStore UserStore, gothicConfig config.Gothic) userHandler {
	oauth.StartOAuth(gothicConfig)
	return userHandler{
		sessionStore: authStore,
		refreshStore: refreshStore,
		userStore:    userStore,
	}
}

func (h userHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if cookie, err := c.Cookie("session"); err == nil {
			s, err := h.sessionStore.ReadAndRefresh(cookie.Value, c.Request().Context())
			if err == nil {
				c.Set("user", s.User)
				return next(c)
			}
			if !errors.Is(err, stores.ErrNotFound) {
				c.Echo().Logger.Error("Failed read session token ", err)
			}
		}

		cookie, err := c.Cookie("refresh")
		if err != nil {
			return next(c)
		}

		refresh, err := h.refreshStore.Read(cookie.Value, c.Request().Context())
		if err != nil {
			if !errors.Is(err, stores.ErrNotFound) {
				return err
			}
			return next(c)
		}
		c.Set("user", refresh.User)

		ses := session.New(uuid.NewString(), refresh.User)
		if err := h.sessionStore.Create(*ses, c.Request().Context()); err != nil {
			return next(c)
		}

		c.SetCookie(ses.Cookie())

		return next(c)
	}
}

func (h userHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/auth/:provider", h.getProvider)
	g.GET("/auth/:provider/callback", h.getCallback)
	g.GET("/me", h.getMe)
	g.POST("/logout", h.logout)
}

func (h userHandler) getProvider(c echo.Context) error {
	ctx := context.WithValue(context.Background(), gothic.ProviderParamKey, c.Param("provider"))

	params := c.QueryParams()
	nonce := base64.URLEncoding.EncodeToString([]byte(rand.Text()))
	params.Set("nonce", nonce)
	state := url.QueryEscape(params.Encode())

	req := c.Request().Clone(ctx)
	req.URL.RawQuery = url.Values{"state": []string{state}}.Encode()

	gothic.BeginAuthHandler(c.Response(), req)

	return nil
}

func (h userHandler) getCallback(c echo.Context) error {
	decodedState, err := url.QueryUnescape(c.QueryParam("state"))
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(decodedState)
	if err != nil {
		return err
	}
	redirectURL := values.Get("redirect")

	gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}

	if strings.ContainsRune(gothUser.NickName, ' ') {
		gothUser.NickName = ""
	}

	u := user.New(gothUser.Email, gothUser.NickName, gothUser.AvatarURL)
	if _, err = h.userStore.ReadOrCreate(u, c.Request().Context()); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}

	ses := session.New(uuid.NewString(), *u)
	if err := h.sessionStore.Create(*ses, c.Request().Context()); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}
	c.SetCookie(ses.Cookie())

	keep, err := strconv.ParseBool(values.Get("keep"))
	if err != nil {
		return err
	}

	ref := refresh.New(uuid.NewString(), *u, keep)
	if err := h.refreshStore.Create(*ref, c.Request().Context()); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}
	c.SetCookie(ref.Cookie())

	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (h userHandler) getMe(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "no user found")
	}

	return c.JSON(http.StatusOK, user)
}

func (h userHandler) logout(c echo.Context) error {
	var errs error
	if refreshCookie, err := c.Cookie("refresh"); err == nil {
		if err := h.refreshStore.Delete(refreshCookie.Value, c.Request().Context()); err != nil {
			c.Echo().Logger.Error("Failed to delete refresh token", err)
			errs = errors.Join(errs, err)
		} else {
			c.SetCookie(&http.Cookie{
				Name:     "refresh",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
		}
	}

	if sessionCookie, err := c.Cookie("session"); err == nil {
		if err := h.sessionStore.Delete(sessionCookie.Value, c.Request().Context()); err != nil {
			c.Echo().Logger.Error("Failed to delete session token", err)
			errs = errors.Join(errs, err)
		} else {
			c.SetCookie(&http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
		}
	}

	if errs != nil {
		return errs
	}

	return c.NoContent(http.StatusOK)
}
