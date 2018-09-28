package version

import (
	"fmt"
)

var (
	BuildTime = "<UNDEFINED>"
	GitCommit = "<UNDEFINED>"
)

func ShowVersion() {
	fmt.Printf("BuildTime: %s.\nGitCommit: %s.\n", BuildTime, GitCommit)
}
