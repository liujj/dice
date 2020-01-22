package global

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/middleware"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var LoggerCron *cron.Cron

func InitLogger() {
	err := os.MkdirAll(LogDir, 0755)
	if err != nil {
		logrus.Warnf("cannot create log dir at %s, err: %v", LogDir, err)
	}

	if !Debug {
		movedLog := ArchiveLog()
		if movedLog > 1 {
			logrus.Infof("archived %d logs from %s to %s", movedLog, LogDir, filepath.Join(LogDir, "archive"))
		}
	}

	logrus.StandardLogger().SetNoLock()

	if !Debug {
		cycleLogger("web")
		cycleLogger("control")
		LoggerCron = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
		LoggerCron.AddFunc("0 0 0 * * *", MidNightLoggerChange)
		LoggerCron.Start()
	} else {
		logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	}
}

func MidNightLoggerChange() {
	cycleLogger("system")
	cycleEchoLogger()
}

func cycleLogger(module string) {
	logName := fmt.Sprintf("%s.%s.%d.log", module, time.Now().Local().Format("060102"), os.Getpid())
	logPath := filepath.Join(LogDir, logName)
	logFd, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil || logFd == nil {
		logrus.Warnf("cannot create stdout file at %s, err: %v", logPath, err)
		return
	}
	var tempLogger *logrus.Logger = nil
	switch module {
	case "system":
		tempLogger = logrus.StandardLogger()
	default:
		logFd.Close()
		return
	}
	oldFd, ok := tempLogger.Out.(*os.File)
	if ok && oldFd != nil {
		if oldFd.Name() != "/dev/stdout" && oldFd.Name() != "/dev/stderr" {
			oldFd.Write([]byte("file closing due to log cycle"))
			oldFd.Close()
		}
	}
	tempLogger.SetOutput(logFd)
}

func cycleEchoLogger() {
	logName := fmt.Sprintf("%s.%s.%d.log", "echo", time.Now().Local().Format("060102"), os.Getpid())
	logPath := filepath.Join(LogDir, logName)
	logFd, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil || logFd == nil {
		logrus.Warnf("cannot create stdout file at %s, err: %v", logPath, err)
		return
	}
	oldFd, ok := middleware.DefaultLoggerConfig.Output.(*os.File)
	middleware.DefaultLoggerConfig.Output = logFd
	if ok && oldFd != nil {
		if oldFd.Name() != "/dev/stdout" && oldFd.Name() != "/dev/stderr" {
			oldFd.Write([]byte("file closing due to log cycle"))
			oldFd.Close()
		}
	}
}

func ArchiveLog() int {
	var i int
	archiveDir := filepath.Join(LogDir, "archive")
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		logrus.Warnf("cannot create log archive dir at %s, err: %v", archiveDir, err)
	}
	filepath.Walk(filepath.Join(LogDir), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			os.Rename(path, filepath.Join(archiveDir, info.Name()))
			i++
		}
		return nil
	})
	return i
}
