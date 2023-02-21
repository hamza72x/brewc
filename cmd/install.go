package cmd

import (
	"fmt"

	"github.com/hamza72x/brewc/pkg/brewc"
	col "github.com/hamza72x/go-color"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a formula",
	Example: `brewc install ffmpeg # for single formulae
brewc install ffmpeg git wget curl # for multiple formulae`,
	Args: cobra.MinimumNArgs(1),
	Run:  runInstallCmd,
}

func init() {
	installCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 10, "number of threads to use for downloading the formulae")
	installCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(installCmd)
}

// runInstallCmd executes the install command.
// Example: brewc install ffmpeg
func runInstallCmd(cmd *cobra.Command, args []string) {

	brewc := brewc.New(_args)

	for _, name := range args {
		fmt.Println("<<<<<<<<<<<< installing", col.Magenta(name), " >>>>>>>>>>>>")
		err := brewc.InstallFormula(name)

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
