package handler

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"net/http"
	"tigerMachine/g"
	"tigerMachine/utils"
)

func HandlePublic(e *echo.Echo) {
	e.GET("/start", Start)
	e.GET("/roll", Roll)
	e.GET("/status", Status)
	e.GET("/lastRecord", LastRecord)
	e.Static("/static", utils.GetRunPath()+g.Config().View)
}
func HandleAdmin(e *echo.Echo) {
	e.GET("/", AppVersion)
	e.GET("/auth", AdminAuth)
	e.GET("/reload", Reload)
	e.GET("/log", ReadLog)
	e.GET("/log/db/:debug", DbLog)
	e.GET("/log/level/:level", LogLevel)
	e.GET("/monitor", ServerMonitor)
	e.POST("/sql", SqlExec)
	e.GET("/key/apply", NewKey)
	e.GET("/key/list", ListKeys)
	e.GET("/key/:key", QueryKey)
	e.GET("/key/:key/set", SetKey)
	e.GET("/key/:key/del", DelKey)
	e.GET("/key/cache/list", QueryKeybyCache)
	e.GET("/key/cache/:key", QueryKeybyCache)
	e.GET("/key/cache/reset", ResetKeyCache)
}

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			msg = fmt.Sprintf("%v, %v", err, he.Internal)
		}
	} else if g.Config().Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	if _, ok := msg.(string); ok {
		msg = map[string]interface{}{
			"httpCode": code,
			"httpMsg":  msg,
			"httpPath": c.Request().RequestURI,
		}
	}
	c.Echo().Logger.Error(err)
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			logrus.Error(err)
		}
	}
}
