package server

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   c.IsTLS() || c.Request().Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}

	return sess
}
