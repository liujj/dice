package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tigerMachine/launcher/cmd"
	"tigerMachine/version"
)

var versionFlag bool

var RootCmd = &cobra.Command{
	Use: "launcher",
	RunE: func(c *cobra.Command, args []string) error {
		if versionFlag {
			version.ShowVersion()
			return nil
		}
		return c.Usage()
	},
}

func init() {
	RootCmd.AddCommand(cmd.Start)
	RootCmd.AddCommand(cmd.Stop)
	RootCmd.AddCommand(cmd.Restart)
	RootCmd.AddCommand(cmd.Check)
	RootCmd.AddCommand(cmd.Monitor)
	RootCmd.AddCommand(cmd.Reload)
	RootCmd.AddCommand(cmd.Auto)

	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "show version")
	cmd.Start.Flags().BoolVar(&cmd.PreqOrderFlag, "preq-order", false, "start modules in the order of prerequisites")
	cmd.Start.Flags().BoolVar(&cmd.ConsoleOutputFlag, "console-output", false, "print the module's output to the console")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
