package cmd

import (
	"flag"
	"fmt"
	"os"

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
	validateInstallCmdArgs()

	brewc := brewc.New(_args.githubToken, _args.threads)

	for _, name := range args {
		fmt.Println("<<<<<<<<<<<< installing", col.Magenta(name), " >>>>>>>>>>>>")
		err := brewc.InstallFormula(name)

		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

// validateInstallCmdArgs parses the command line flags and returns a cmdArg struct.
func validateInstallCmdArgs() {

	if len(_args.githubToken) == 0 {
		_args.githubToken = os.Getenv("HOMEBREW_GITHUB_API_TOKEN")
	}

	if len(_args.githubToken) == 0 {
		fmt.Println("Either set the HOMEBREW_GITHUB_API_TOKEN env variable or pass the github-token flag")
		flag.Usage()
		os.Exit(1)
	}
}
