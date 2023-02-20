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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a formula",
	Example: `brewc install ffmpeg # for single formulae
brewc install ffmpeg git wget curl # for multiple formulae`,
	Args: cobra.MinimumNArgs(1),
	Run:  runInstallCmd,
}

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
	installCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 10, "number of threads to use for downloading the formulae")
	installCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
}

// Run executes the root command.
func Run() {
	constant.Initialize()

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
