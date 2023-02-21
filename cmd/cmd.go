package cmd

import (
	"os"

	"github.com/hamza72x/brewc/pkg/constant"
	"github.com/hamza72x/brewc/pkg/models"
	"github.com/spf13/cobra"
)

// _args holds the optional arguments passed to the command line.
var _args = &models.OptionalArgs{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "brewc",
	Short: "Install brew packages with concurrent downloads instead of one by one (which is typically slow)",
	Run:   runRootCmd,
}

// Run executes the root command.
func Run() {
	archAndCodeName := models.GetArchAndOSName()

	constant.Initialize(archAndCodeName.Architecture)

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
// only used to print the usage.
func runRootCmd(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
