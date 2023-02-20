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
	// install cmd
	installCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 5, "number of threads to use for downloading the formulae")
	installCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(installCmd)

	// uninstall cmd
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 5, "number of threads to use for downloading the formulae")
	uninstallCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")
	uninstallCmd.Flags().BoolVarP(&_args.DeleteUnusedDependencies, "delete-unused-dependencies", "d", false, "delete unused dependencies after uninstalling a formula")

	// reinstall cmd
	rootCmd.AddCommand(reinstallCmd)

	reinstallCmd.Flags().IntVarP(&_args.Threads, "threads", "t", 5, "number of threads to use for downloading the formulae")
	reinstallCmd.Flags().BoolVarP(&_args.Verbose, "verbose", "v", false, "verbose output")
}

// Run executes the root command.
func Run() {
	constant.Initialize()

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
