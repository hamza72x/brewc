package brewc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hamza72x/brewc/pkg/brew"
	"github.com/hamza72x/brewc/pkg/models"
	"github.com/hamza72x/brewc/pkg/models/formula"
	col "github.com/hamza72x/go-color"
)

var (
	greenArrow = col.Green("==>")
	redArrow   = col.Red("==>")
)

// BrewC downloads all of the dependencies for a formula in concurrent goroutines.
// and then handles to `brew` command to install the formula.
type BrewC struct {
	// threads is the number of concurrent goroutines used to download the formula dependencies.
	threads int

	archAndCodeName *models.ArchAndCodeName

	// httpClient is the http client used to make requests to the GitHub API.
	// usually to calculate the download URL for a formula.
	httpClient *http.Client

	// brew is the brew command wrapper
	brew *brew.Brew

	args *models.OptionalArgs
}

// New returns a new BrewC instance.
func New(args *models.OptionalArgs) *BrewC {
	archAndCodeName := models.GetArchAndOSName()

	return &BrewC{
		threads:         args.Threads,
		archAndCodeName: archAndCodeName,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		brew: brew.New(),
		args: args,
	}
}

// InstallFormula installs the given formula.
// Example: InstallFormula("ffmpeg")
func (b *BrewC) InstallFormula(name string) error {

	_, err := b.FormulasReverseIterate(name, false, func(index int, f *formula.Formula) {
		if f.Name == name {
			return
		}
		if err := b.brew.InstallFormula(f.Name, b.args.Verbose); err != nil {
			fmt.Printf("%s Error installing formula: %s\n", redArrow, err.Error())
		}
	})

	if err != nil {
		return err
	}

	// last formula is the formula we want to install
	if err := b.brew.InstallFormula(name, b.args.Verbose); err != nil {
		return err
	}

	// install the formula by calling the `brew` command
	// return b.brew.InstallFormula(name, b.verbose)
	return nil
}

// UninstallFormula uninstalls the given formula.
// Example: UninstallFormula("ffmpeg")
func (b *BrewC) UninstallFormula(name string) error {

	if err := b.brew.UninstallFormula(name, b.args.Verbose); err != nil {
		fmt.Printf("%s Error uninstalling formula: %s\n", redArrow, err.Error())
	}

	_, err := b.FormulasReverseIterate(name, true, func(index int, f *formula.Formula) {
		if err := b.brew.UninstallFormula(f.Name, b.args.Verbose); err != nil {
			fmt.Printf("%s Error uninstalling formula: %s\n", redArrow, err.Error())
		}
	})

	if err != nil {
		return err
	}

	return nil
}
