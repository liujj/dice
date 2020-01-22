package version

import (
	"fmt"
	"runtime"
)

var (
	LdFlagUndefined = "<UNDEFINED>"
	BuildTime       = LdFlagUndefined
	GitCommit       = LdFlagUndefined
	LatestTag       = LdFlagUndefined
)

func ShowVersion() {
	fmt.Printf("BuildTime: %s\nGitCommit: %s\nLatestTag: %s\nGoVersion: %s\n", BuildTime, GitCommit, LatestTag, runtime.Version())
}
