package brewc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hamza72x/brewc/pkg/brew"
	"github.com/hamza72x/brewc/pkg/constant"
	"github.com/hamza72x/brewc/pkg/models"
	"github.com/hamza72x/brewc/pkg/models/formula"
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

	list, err := formula.GetFormulaList(name, &formula.GetFormulaListOpts{
		IncludeInstalled: false,
		DependencyLevel:  -1,
		Threads:          b.threads,
	})

	if err != nil {
		return err
	}

	fmt.Println("")

	list.IterateChildFirst(b.threads, func(f *formula.Formula) {

		fmt.Printf("%s Working On: %s\n", constant.GreenArrow, f.Name)

		if err := b.brew.InstallFormula(f.Name, b.args.Verbose); err != nil {
			fmt.Printf("%s Error installing formula (%s): %s\n", constant.RedArrow, f.Name, err.Error())
		}
	})

	return nil
}

// UninstallFormula uninstalls the given formula.
// Example: UninstallFormula("ffmpeg")
func (b *BrewC) UninstallFormula(name string) error {

	if !b.args.DeleteUnusedDependencies && !b.args.DeleteAllNestedDependencies {
		return b.brew.UninstallFormula(name, b.args.Verbose)
	}

	dependencyLevel := 1

	if b.args.DeleteAllNestedDependencies {
		dependencyLevel = -1
	}

	list, err := formula.GetFormulaList(name, &formula.GetFormulaListOpts{
		IncludeInstalled: true,
		DependencyLevel:  dependencyLevel,
		Threads:          b.threads,
		Unique:           true,
	})

	if err != nil {
		return err
	}

	fmt.Println("")

	list.IterateParentFirst(b.threads, func(f *formula.Formula) {
		fmt.Printf("%s Removing: %s\n", constant.GreenArrow, f.Name)
		if err := b.brew.UninstallFormula(f.Name, b.args.Verbose); err != nil && b.args.Verbose {
			fmt.Printf("%s Error uninstalling formula: %s\n", constant.RedArrow, err.Error())
		}
	})

	return nil
}

// ReinstallFormula uninstalls and then installs the given formula.
// Example: ReinstallFormula("ffmpeg")
func (b *BrewC) ReinstallFormula(name string) error {
	return b.brew.ReinstallFormula(name, b.args.Verbose)
}
