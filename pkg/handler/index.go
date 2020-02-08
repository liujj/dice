package handler

import (
	"dice/pkg/global"
	_ "dice/pkg/statik"
	"dice/pkg/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
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
	var disableEmbedResource bool
	if global.Debug {
		disableEmbedResource = true
		possiblePath1 := filepath.Join(utils.GetRunPath(), "./view")
		possiblePath2 := filepath.Join(utils.GetRunPath(), "../view")
		if _, err := os.Stat(possiblePath1); err == nil {
			e.Static("/static", possiblePath1)
		} else if _, err := os.Stat(possiblePath2); err == nil {
			e.Static("/static", possiblePath2)
		} else {
			disableEmbedResource = false
		}
	}
	if !disableEmbedResource {
		statikFS, err := fs.New()
		if err != nil {
			logrus.Fatalf("Static Resources load error: %v", err)
		}
		e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(statikFS))))
	}
}

func Root(c echo.Context) error {
	return c.Redirect(http.StatusFound, "./static/index.html")
}
