package command

import (
	"dice/pkg/global"
	"dice/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

var normalArgs []string
var myName = filepath.Base(os.Args[0])

func Parse() {
	var cmd []string
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if utils.InStringArray(arg, []string{"init", "start", "status", "stop", "restart", "remove"}) {
			cmd = append(cmd, arg)
		} else {
			normalArgs = append(normalArgs, arg)
		}
	}
	if len(cmd) == 0 {
		return
	}
	if len(cmd) > 1 {
		logrus.Errorf("detect two command: %v", cmd)
		os.Exit(1)
	}
	switch cmd[0] {
	case "start":
		Start()
	case "status":
		Status(nil)
	case "stop":
		Stop()
	case "restart":
		Restart()
	case "remove":
		Remove()
	}
	os.Exit(0)
}

func writePid(pid int) {
	_ = os.MkdirAll(filepath.Dir(global.PidFile), 0755)
	fi, err := os.OpenFile(global.PidFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err == nil {
		defer fi.Close()
		_, _ = fi.WriteString(strconv.Itoa(pid))
	}
}

func deletePid() {
	err := os.Remove(global.PidFile)
	if err != nil {
		logrus.Errorf("error removing pid file, error: %v", err)
	}
}

func fetchPid() int {
	var pidRaw []byte
	pidRaw, err := ioutil.ReadFile(global.PidFile)
	if err != nil {
		return -1
	}
	pid, err := strconv.Atoi(string(pidRaw))
	if err != nil {
		return -1
	}
	return pid
}

type blackhole struct{}

func (b blackhole) Write(p []byte) (n int, err error) {
	return len(p), nil
}
