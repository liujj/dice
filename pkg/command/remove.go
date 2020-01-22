package command

import (
	"dice/pkg/global"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func Remove() {
	noOutputLogger := logrus.New()
	noOutputLogger.Out = blackhole{}
	runningPid := Status(noOutputLogger)
	if runningPid != -1 {
		Stop()
	}

	abs, absErr := filepath.Abs(os.Args[0])
	if absErr != nil {
		logrus.Errorf("Cannot")
		return
	}
	filepath.Walk(filepath.Dir(abs), func(path string, info os.FileInfo, err error) error {
		if filepath.Dir(abs) == path {
			return nil
		}
		os.RemoveAll(path)
		return nil
	})
	logrus.Infof("removed")
	logrus.Infof("don't forget to remove log files in %s", global.LogDir)
}
