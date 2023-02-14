package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type optionalArgs struct {
	githubToken string
	threads     int
}

// _args holds the optional arguments passed to the command line.
var _args = optionalArgs{}

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

func init() {
	installCmd.Flags().StringVarP(&_args.githubToken, "github-token", "g", "", "github token to use for downloading the formulae")
	installCmd.Flags().IntVarP(&_args.threads, "threads", "t", 2, "number of threads to use for downloading the formulae")
	rootCmd.AddCommand(installCmd)
}

// Run executes the root command.
func Run() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
