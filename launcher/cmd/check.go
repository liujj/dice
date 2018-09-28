package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tigerMachine/launcher/g"
)

var Check = &cobra.Command{
	Use:   "check [Module ...]",
	Short: "Check the status of modules",
	Long:  "Check if the specified modules are running.",
	RunE:  check,
}

func check(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			return fmt.Errorf("%s doesn't exist", moduleName)
		}

		if g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
		} else {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", "DOWN", "\n")
		}
	}

	return nil
}
