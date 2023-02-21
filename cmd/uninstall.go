package cmd

import (
	"fmt"

	"github.com/hamza72x/brewc/pkg/brewc"
	col "github.com/hamza72x/go-color"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall a formula",
	Example: `brewc uninstall ffmpeg # for single formulae
brewc uninstall ffmpeg git wget curl # for multiple formulae`,
	Args: cobra.MinimumNArgs(1),
	Run:  runUninstallCmd,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 10, "number of threads to use for downloading the formulae")
	uninstallCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")
	uninstallCmd.Flags().BoolVarP(&_args.DeleteUnusedDependencies, "delete-unused-dependencies", "d", false, "delete unused dependencies after uninstalling a formula")
	uninstallCmd.Flags().BoolVarP(&_args.DeleteAllNestedDependencies, "delete-all-nested-dependencies", "D", false, "delete all sub-dependencies after uninstalling a formula (it will delete all nested unused dependencies)")
}

// runUninstallCmd executes the uninstall command.
// Example: brewc uninstall ffmpeg
func runUninstallCmd(cmd *cobra.Command, args []string) {

	brewc := brewc.New(_args)

	for _, name := range args {
		fmt.Println("<<<<<<<<<<<< uninstalling", col.Magenta(name), " >>>>>>>>>>>>")
		err := brewc.UninstallFormula(name)

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
