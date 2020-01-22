package handler

import (
	"dice/pkg/global"
	"dice/pkg/utils"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func HandlePublic(e *echo.Group) {
	e.GET("/start", Start)
	e.GET("/roll", Roll)
	e.GET("/status", Status)
	e.GET("/lastRecord", LastRecord)
}
func HandleAdmin(e *echo.Group) {
	e.GET("/auth", AdminAuth)
	e.GET("/reload", Reload)
	e.GET("/log", ReadLog)
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
func HandleGlobal(e *echo.Group) {
	e.GET("/", Root)
	e.GET("/version", AppVersion)
	e.Static("/static", filepath.Join(utils.GetRunPath(), global.Config().View))
}
func Root(c echo.Context) error {
	return c.Redirect(http.StatusFound, "./static/index.html")
}
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}
