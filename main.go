package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	echoLog "github.com/labstack/gommon/log"
	"net/http"
	"os"
	"tigerMachine/g"
	"tigerMachine/handler"
	"tigerMachine/model"
	"tigerMachine/version"
)

var (
	ePublic *echo.Echo
	eAdmin  *echo.Echo
)

func main() {
	processFlag()
	initEcho()
	model.InitDatabase()
	handler.HandleAdmin(eAdmin)
	handler.HandlePublic(ePublic)
	logrus.Infof("Admin http server started on %s", g.Config().Http.AdminPort)
	logrus.Infof("Public http server started on %s", g.Config().Http.PublicPort)
	launchErr := gracehttp.Serve(
		&http.Server{Addr: g.Config().Http.PublicPort, Handler: ePublic},
		&http.Server{Addr: g.Config().Http.AdminPort, Handler: eAdmin},
	)
	if launchErr != nil {
		logrus.Fatalln(launchErr)
	}
}

func initEcho() {
	//timeStr := time.Now().Local().Format(model.TimeFormartNums)
	accessLogFile, ferr := os.OpenFile(g.Config().Log+"access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if ferr == nil {
		middleware.DefaultLoggerConfig.Output = accessLogFile
	} else {
		fmt.Println(ferr.Error())
	}
	ePublic = echo.New()
	ePublic.Use(middleware.Recover())
	ePublic.Use(middleware.Logger())
	if g.Config().Debug {
		ePublic.Logger.SetLevel(echoLog.DEBUG)
	} else {
		ePublic.Logger.SetLevel(echoLog.WARN)
	}
	ePublic.Use(session.Middleware(sessions.NewCookieStore([]byte("lalala"))))
	ePublic.Use(model.StatMiddleWare)
	ePublic.Use(handler.PublicAuthMiddleWare())
	ePublic.HTTPErrorHandler = handler.CustomHTTPErrorHandler
	if g.Config().Http.GzipLevel > 0 {
		ePublic.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: g.Config().Http.GzipLevel,
		}))
	}

	eAdmin = echo.New()
	eAdmin.Use(middleware.Recover())
	eAdmin.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	eAdmin.Use(handler.AdminAuthMiddleWare())
}

func processFlag() {
	model.InitStat()
	configPath := flag.String("c", "./config.json", "Configuration file path")
	showVersion := flag.Bool("v", false, "Show version")
	neepHelp := flag.Bool("h", false, "Show this help")
	flag.Parse()
	if *showVersion {
		version.ShowVersion()
		os.Exit(0)
	}
	if *neepHelp {
		flag.Usage()
		os.Exit(0)
	}
	g.ParseConfig(*configPath)
}
