package main

import (
	"dice/pkg/command"
	"dice/pkg/global"
	"dice/pkg/handler"
	"dice/pkg/model"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/sirupsen/logrus"
)

func main() {
	global.ParseFlag()
	command.Parse()
	global.ParseConfig()
	global.InitLogger()
	handler.InitEcho()
	model.InitDatabase()
	model.InitStat()
	logrus.Infof("http server started on %s", global.Config().Http.Port)
	launchErr := gracehttp.Serve(
		&http.Server{Addr: global.Config().Http.Port, Handler: handler.E},
	)
	if launchErr != nil {
		logrus.Fatalln(launchErr)
	}
}
