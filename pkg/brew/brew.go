package brew

import (
	"fmt"

	"github.com/hamza72x/brewc/pkg/util"
	col "github.com/hamza72x/go-color"
)

type Brew struct {
	bin string
}

func New() *Brew {
	return &Brew{
		bin: getBrewBinary(),
	}
}

// InstallFormula installs the given formula.
func (b *Brew) InstallFormula(name string, verbose bool) error {
	var args = []string{"install", name}

	if verbose {
		args = append(args, "-v")
	}

	return util.ExecStandard(b.bin, args...)
}

// getBrewBinary returns the path to the brew binary.
func getBrewBinary() string {
	var bin string

	var paths = []string{
		"/usr/local/bin/brew",
		"/opt/homebrew/bin/brew",
		"/home/linuxbrew/.linuxbrew/bin/brew",
	}

	for _, path := range paths {
		if util.DoesFileExist(path) {
			bin = path
			break
		}
	}

	if len(bin) == 0 {
		fmt.Printf("brew binary not found in any of the following paths: %+v", paths)
		// FIXME: exit with error
		// os.Exit(1)
	}

	fmt.Printf("%s: %s\n", col.Green("Brew Binary"), bin)

	return bin
}
