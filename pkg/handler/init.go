package handler

import (
	"dice/pkg/global"
	"dice/pkg/model"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

var (
	E       *echo.Echo
	eGlobal *echo.Group
	ePublic *echo.Group
	eAdmin  *echo.Group
)

func InitEcho() {
	E = echo.New()
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(E)
	E.Use(middleware.Recover())
	E.Use(middleware.Logger())
	if global.Debug {
		E.Logger.SetLevel(echoLog.DEBUG)
	} else {
		E.Logger.SetLevel(echoLog.WARN)
	}
	eGlobal = E.Group("")
	eGlobal.Use(middleware.Recover())
	ePublic = E.Group("api")
	ePublic.Use(session.Middleware(sessions.NewCookieStore([]byte("lalala"))))
	ePublic.Use(model.StatMiddleWare)
	ePublic.Use(PublicAuthMiddleWare())
	if global.Config().Http.GzipLevel > 0 {
		ePublic.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: global.Config().Http.GzipLevel,
		}))
	}
	eAdmin = E.Group("admin")
	eAdmin.Use(middleware.Recover())
	eAdmin.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	eAdmin.Use(AdminAuthMiddleWare())

	HandleGlobal(eGlobal)
	HandlePublic(ePublic)
	HandleAdmin(eAdmin)
}
