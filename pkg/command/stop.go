package command

import (
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

func Stop() {
	noOutputLogger := logrus.New()
	noOutputLogger.Out = blackhole{}
	runningPid := Status(noOutputLogger)
	if runningPid == -1 {
		logrus.Warnf("not running, run `status` to show detail")
		return
	}

	cmd := exec.Command("kill", strconv.Itoa(runningPid))
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		logrus.Errorf("cannot execute kill, error: %v", err)
		return
	}
	state, err := cmd.Process.Wait()
	if err != nil || state == nil {
		logrus.Errorf("error executing kill, error: %v", err)
		return
	}
	if state.Success() {
		logrus.Infof("stopping pid %d", runningPid)
	} else {
		logrus.Errorf("error stopping pid %d, exitcode: %v", runningPid, state.ExitCode())
		return
	}
	runningPid = Status(noOutputLogger)
	if runningPid == -1 {
		logrus.Infof("stopped")
		deletePid()
		return
	} else {
		logrus.Warnf("pid %d is still running", runningPid)
		return
	}
}
