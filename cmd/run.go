package cmd

import (
	"fmt"

	"github.com/hamza72x/brewc/pkg/brewc"
	col "github.com/hamza72x/go-color"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// only used to print the usage.
func runRootCmd(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// runInstallCmd executes the install command.
// Example: brewc install ffmpeg
func runInstallCmd(cmd *cobra.Command, args []string) {

	brewc := brewc.New(_args)

	for _, name := range args {
		fmt.Println("<<<<<<<<<<<< installing", col.Magenta(name), " >>>>>>>>>>>>")
		err := brewc.InstallFormula(name)

		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

// runUninstallCmd executes the uninstall command.
// Example: brewc uninstall ffmpeg
func runUninstallCmd(cmd *cobra.Command, args []string) {

	brewc := brewc.New(_args)

	for _, name := range args {
		fmt.Println("<<<<<<<<<<<< uninstalling", col.Magenta(name), " >>>>>>>>>>>>")
		err := brewc.UninstallFormula(name)

		if err != nil {
			fmt.Println("error:", err)
		}
	}
}
