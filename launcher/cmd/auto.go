package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"tigerMachine/launcher/g"
	"time"
)

var Auto = &cobra.Command{
	Use:   "auto [Module ...]",
	Short: "Start or Reload modules",
	Long:  "Start if stopped, reload if running.",
	RunE:  auto,
}

func auto(c *cobra.Command, args []string) error {
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
			start(c, []string{moduleName})
			continue
		}
		check(c, []string{moduleName})
		cmd := exec.Command("kill", "-USR2", g.Pid(moduleName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err == nil {
			fmt.Print("[", g.ModuleApps[moduleName], "] reloading\n")
			time.Sleep(2 * time.Second)
			check(c, []string{moduleName})
			continue
		}
		return err
	}
	return nil
}
