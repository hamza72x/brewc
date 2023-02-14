package cmd

import "github.com/hamza72x/brewc/pkg/brewc"

func Run() {
	brewC := brewc.New("")

	brewC.InstallFormula("git")
}
