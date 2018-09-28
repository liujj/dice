package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"tigerMachine/launcher/g"
)

var Stop = &cobra.Command{
	Use:   "stop [Module ...]",
	Short: "Stop modules",
	Long:  "Stop the specified modules.",
	RunE:  stop,
}

func stop(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			return fmt.Errorf("%s doesn't exist", moduleName)
		}

		if !g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] not running\n")
			continue
		}

		cmd := exec.Command("kill", "-TERM", g.Pid(moduleName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err == nil {
			fmt.Print("[", g.ModuleApps[moduleName], "] DOWN("+g.Pid(moduleName)+")\n")
			continue
		}
		return err
	}
	return nil
}
