package command

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func Status(logger *logrus.Logger) int {
	if logger == nil {
		logger = logrus.New()
	}

	pid := fetchPid()
	if pid == -1 {
		logger.Warn("cannot find pid file")
		return -1
	}
	out := bytes.NewBuffer([]byte{})
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
	cmd.Stdin = nil
	cmd.Stdout = out
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		logger.Infof("found pid from file: %d", pid)
		logger.Error("cannot execute ps, error: %v", err)
		return -1
	}
	state, err := cmd.Process.Wait()
	if err != nil || state == nil {
		logger.Infof("found pid from file: %d", pid)
		logger.Error("error executing ps, error: %v", err)
		return -1
	}
	if state.Success() {
		if !strings.Contains(out.String(), myName) {
			logger.Infof("found pid from file: %d", pid)
			logger.Warnf("process name does not contain '%s'", myName)
			logger.Infof("not running")
			return -1
		}
		logger.Infof("running with pid %d", pid)
		return pid
	} else {
		logger.Infof("found pid from file: %d", pid)
		logger.Infof("not running")
		return -1
	}
}
