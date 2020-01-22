package command

import (
	"dice/pkg/global"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

func Start() {
	noOutputLogger := logrus.New()
	noOutputLogger.Out = blackhole{}
	runningPid := Status(noOutputLogger)
	if runningPid != -1 {
		logrus.Warnf("already running with pid %d, please check", runningPid)
		return
	}

	err := os.MkdirAll(global.LogDir, 0755)
	if err != nil {
		logrus.Warnf("cannot create stdout dir at %s, err: %v", global.LogDir, err)
	}

	stdoutPath := filepath.Join(global.LogDir, "stdout.log")
	stdoutFd, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Warnf("cannot create stdout file at %s, err: %v", stdoutPath, err)
		return
	}
	stderrPath := filepath.Join(global.LogDir, "stderr.log")
	stderrFd, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Warnf("cannot create stderr file at %s, err: %v", stderrPath, err)
		return
	}

	cmd := exec.Command(os.Args[0], normalArgs...)
	cmd.Stdin = nil
	cmd.Stdout = stdoutFd
	cmd.Stderr = stderrFd
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	err = cmd.Start()
	if err != nil {
		logrus.Warnf("failed to start, please run `status` to check")
		logrus.Infof("prosibble pid: %d", cmd.Process.Pid)
	}
	writePid(cmd.Process.Pid)
	logrus.Infof("started with pid %d", cmd.Process.Pid)
	cmd.Process.Release()
}
