package handlers

import (
	"cmp"
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
)

type (
	SessionStore interface {
		Create(ctx context.Context, session session.Session) error
		ReadAndRefresh(ctx context.Context, uuid string) (session session.Session, err error)
		Delete(ctx context.Context, uuid string) error
	}

	RefreshStore interface {
		Create(ctx context.Context, token refresh.Refresh) error
		Read(ctx context.Context, uuid uuid.UUID) (refresh.Refresh, error)
		Delete(ctx context.Context, uuid uuid.UUID) error
	}

	UserStore interface {
		ReadOrCreate(ctx context.Context, user user.User) (u user.User, isNew bool, err error)
	}

	userHandler struct {
		sessionStore SessionStore
		refreshStore RefreshStore
		userStore    UserStore
	}
)

func NewUserHandler(authStore SessionStore, refreshStore RefreshStore, userStore UserStore, gothicConfig config.Gothic) userHandler {
	oauth.StartOAuth(gothicConfig)
	return userHandler{
		sessionStore: authStore,
		refreshStore: refreshStore,
		userStore:    userStore,
	}
}

func (h userHandler) RegisterRoutes(g *echo.Group, protection echo.MiddlewareFunc) {
	g.GET("/auth/:provider", h.getProvider)
	g.GET("/auth/:provider/callback", h.getCallback)
	g.GET("/me", h.getMe, protection)
	g.POST("/logout", h.logout, protection)
}

func (h userHandler) getUserFromSession(c echo.Context) (u user.User, err error) {
	cookie, err := c.Cookie("session")
	if err := cmp.Or(err, cookie.Valid()); err != nil {
		return u, err
	}

	session, err := h.sessionStore.ReadAndRefresh(c.Request().Context(), cookie.Value)
	if err != nil {
		return u, err
	}

	return session.User, nil
}

func (h userHandler) getUserFromRefresh(c echo.Context) (u user.User, err error) {
	cookie, err := c.Cookie("refresh")
	if err := cmp.Or(err, cookie.Valid()); err != nil {
		return u, err
	}

	cookieUUID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return u, err
	}

	refresh, err := h.refreshStore.Read(c.Request().Context(), cookieUUID)
	if err != nil {
		return u, err
	}

	return refresh.User, nil
}

func (h userHandler) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := h.getUserFromSession(c)
		if err == nil {
			c.Set("user", u)
			return next(c)
		}

		u, err = h.getUserFromRefresh(c)
		if err != nil {
			return next(c)
		}
		c.Set("user", u)

		ses := session.New(uuid.NewString(), u)
		if h.sessionStore.Create(c.Request().Context(), *ses) != nil {
			c.SetCookie(ses.Cookie())
		}

		return next(c)
	}
}

func (h userHandler) Protection(next echo.HandlerFunc) echo.HandlerFunc {
	return h.Authentication(
		func(c echo.Context) error {
			if _, ok := c.Get("user").(user.User); !ok {
				return echo.ErrUnauthorized.SetInternal(
					errors.New("Protection: no user in context"))
			}
			return next(c)
		},
	)
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

	tmpUser := *user.New(gothUser.Email, gothUser.NickName, gothUser.AvatarURL)
	u, _, err := h.userStore.ReadOrCreate(c.Request().Context(), tmpUser)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}

	ses := *session.New(uuid.NewString(), u)
	if err := h.sessionStore.Create(c.Request().Context(), ses); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}
	c.SetCookie(ses.Cookie())

	keep, err := strconv.ParseBool(values.Get("keep"))
	if err != nil {
		return err
	}

	ref := *refresh.New(uuid.New(), u, keep)
	if err := h.refreshStore.Create(c.Request().Context(), ref); err != nil {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}
	c.SetCookie(ref.Cookie())

	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (h userHandler) getMe(c echo.Context) error {
	return c.JSON(http.StatusOK, MustGetUser(c))
}

func (h userHandler) logout(c echo.Context) error {
	var errs error

	if refreshCookie, err := c.Cookie("refresh"); err == nil {
		if id, err := uuid.Parse(refreshCookie.Value); err != nil {
			c.Echo().Logger.Error("Invalid refresh token", err)
			errs = errors.Join(errs, err)
		} else if err := h.refreshStore.Delete(c.Request().Context(), id); err != nil {
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
		if err := h.sessionStore.Delete(c.Request().Context(), sessionCookie.Value); err != nil {
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

func MustGetUser(c echo.Context) user.User {
	return c.Get("user").(user.User)
}
