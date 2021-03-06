package handler

import (
	"dice/pkg/global"
	"encoding/base64"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AdminAuthMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("admin_session", c)
			if sess.Values["admin"] == "1" || c.Path() == "/admin/auth" {
				return next(c)
			} else {
				return echo.ErrForbidden
			}
		}
	}
}

func auth(c echo.Context) error {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	l := len("basic")

	if len(auth) > l+1 && strings.ToLower(auth[:l]) == "basic" {
		b, err := base64.StdEncoding.DecodeString(auth[l+1:])
		if err != nil {
			return err
		}
		cred := string(b)
		for i := 0; i < len(cred); i++ {
			if cred[i] == ':' {
				if cred[:i] == global.Config().Admin.User && cred[i+1:] == global.Config().Admin.Pass {
					return nil
				}
				break
			}
		}
	}
	c.Response().Header().Set(echo.HeaderWWWAuthenticate, "basic realm=Restricted")
	return echo.ErrUnauthorized
}

func PublicAuthMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("pubilc_session", c)
			_, hasAk := sess.Values["ak"]
			if hasAk || c.Path() == "/api/start" {
				return next(c)
			} else {
				return echo.ErrForbidden
			}
		}
	}
}
