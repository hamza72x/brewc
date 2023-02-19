package brewc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hamza72x/brewc/pkg/brew"
	"github.com/hamza72x/brewc/pkg/models"
	"github.com/hamza72x/brewc/pkg/models/formula"
	"github.com/hamza72x/brewc/pkg/util"
	col "github.com/hamza72x/go-color"
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

	// downloadedData tracks if certain formula is downloaded or not
	// key: string is the formula name
	downloadedData map[string]bool

	verbose bool
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
		brew:    brew.New(),
		verbose: args.Verbose,
	}
}

// InstallFormula installs the given formula.
// Example: InstallFormula("ffmpeg")
func (b *BrewC) InstallFormula(name string) error {
	// return b.brew.InstallFormula(formula)
	list, err := b.GetAllFormulas(name)

	if err != nil {
		return err
	}

	// download all of the formula & dependencies
	// err = b.DownloadFormulas(list)
	var wg sync.WaitGroup
	var c = make(chan int, b.threads/2)

	list.IterateReverse(func(index int, f *formula.Formula) {
		if f.Name == name {
			return
		}

		wg.Add(1)
		go func(f *formula.Formula) {
			c <- 1
			defer wg.Done()
			defer func() { <-c }()
			if err := b.brew.InstallFormula(f.Name, b.verbose); err != nil {
				fmt.Printf("error installing formula: %s", err.Error())
			}
		}(f)
	})

	wg.Wait()

	// last formula is the formula we want to install
	if err := b.brew.InstallFormula(name, b.verbose); err != nil {
		return err
	}

	// install the formula by calling the `brew` command
	// return b.brew.InstallFormula(name, b.verbose)
	return nil
}

// GetAllFormulas returns all of the formulas.
func (b *BrewC) GetAllFormulas(name string) (*formula.FormulaList, error) {
	var list = formula.NewFormulaList()

	var err = b.setAllFormulasRecursive(name, list, 0)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// setAllFormulasRecursive returns all of the formulas recursively.
func (b *BrewC) setAllFormulasRecursive(name string, list *formula.FormulaList, nestedCount int) error {

	var wg sync.WaitGroup
	var conn = make(chan int, b.threads)

	var mainFormula, err = b.getFormulaJSON(name)

	if err != nil {
		return err
	}

	// if the formula is already installed, then we don't need to install it again.
	// that's also means that all of its dependencies are already installed too.
	if mainFormula.IsInstalled() {
		return nil
	}

	if !list.HasFormula(mainFormula) {
		list.Add(mainFormula)
	}

	for _, dep := range mainFormula.Dependencies {
		wg.Add(1)

		go func(dep string) {
			conn <- 1

			defer wg.Done()
			defer func() { <-conn }()

			f, err := b.getFormulaJSON(dep)

			if err != nil {
				fmt.Printf("error getting formula %s: %v", dep, err)
				return
			}

			if f.IsInstalled() {
				return
			}

			if !list.HasFormula(f) {
				fmt.Printf("%s ", col.Info(dep))
				list.Add(f)
			}

			if err := b.setAllFormulasRecursive(dep, list, nestedCount+1); err != nil {
				fmt.Printf("error getting nested formula %s: %v, %d", dep, err, nestedCount)
				return
			}

		}(dep)
	}

	wg.Wait()

	return nil
}

// DECIDE: should we use the github API to get the list of formulas?
// Or check local installation folder
func (b *BrewC) getFormulaJSON(name string) (*formula.Formula, error) {
	var f formula.Formula

	var url = util.GetFormulaURL(name)

	resp, err := b.doGET(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&f)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

// doGET makes a GET request to the given URL.
func (b *BrewC) doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
