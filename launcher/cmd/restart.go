package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"tigerMachine/launcher/g"
)

var Restart = &cobra.Command{
	Use:   "restart [Module ...]",
	Short: "Restart modules",
	Long:  "Restart the specified or all modules",
	RunE:  restart,
}

func restart(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if err := stop(c, []string{moduleName}); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		if err := start(c, []string{moduleName}); err != nil {
			return err
		}
	}
	return nil
}
