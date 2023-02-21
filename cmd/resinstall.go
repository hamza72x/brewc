package cmd

import (
	"fmt"

	"github.com/hamza72x/brewc/pkg/brewc"
	col "github.com/hamza72x/go-color"
	"github.com/spf13/cobra"
)

// reinstallCmd represents the reinstall command
var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "reinstall a formula",
	Example: `brewc reinstall ffmpeg # for single formulae
brewc reinstall ffmpeg git wget curl # for multiple formulae`,
	Args: cobra.MinimumNArgs(1),
	Run:  runReinstallCmd,
}

func init() {

	rootCmd.AddCommand(reinstallCmd)

	reinstallCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 5, "number of threads to use for downloading the formulae")
	reinstallCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")
}

// runReinstallCmd executes the reinstall command.
// Example: brewc reinstall ffmpeg
func runReinstallCmd(cmd *cobra.Command, args []string) {

	brewc := brewc.New(_args)

	for _, name := range args {
		fmt.Println("<<<<<<<<<<<< reinstalling", col.Magenta(name), " >>>>>>>>>>>>")
		err := brewc.ReinstallFormula(name)

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
