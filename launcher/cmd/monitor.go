package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"tigerMachine/launcher/g"
)

var Monitor = &cobra.Command{
	Use:   "monitor [Module ...]",
	Short: "Display module's log",
	Long:  "Display a specified module's log.",
	RunE:  monitor,
}

func checkMonReq(name string) error {
	if !g.HasModule(name) {
		return fmt.Errorf("%s doesn't exist", name)
	}

	if !g.HasLogfile(name) {
		r := g.Rel(g.Cfg(name))
		return fmt.Errorf("expect logfile: %s", r)
	}

	return nil
}

func monitor(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		if len(g.AllModulesInOrder) == 1 {
			args = g.AllModulesInOrder
		} else {
			return c.Usage()
		}
	}
	var tailArgs []string = []string{"-f"}
	for _, moduleName := range args {
		if err := checkMonReq(moduleName); err != nil {
			return err
		}

		tailArgs = append(tailArgs, g.LogPath(moduleName))
	}
	cmd := exec.Command("tail", tailArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
