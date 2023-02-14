package brewc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hamza72x/brewc/pkg/brew"
	"github.com/hamza72x/brewc/pkg/downloader"
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

	archOSName *models.ArchAndOS

	// githubToken is the token used to authenticate with the GitHub API.
	// required permissions: read:packages
	githubToken string

	// httpClient is the http client used to make requests to the GitHub API.
	// usually to calculate the download URL for a formula.
	httpClient *http.Client

	// brew is the brew command wrapper
	brew *brew.Brew

	// downloader is the downloader used to download the formula dependencies.
	downloader *downloader.Downloader

	// downloadedData tracks if certain formula is downloaded or not
	// key: string is the formula name
	downloadedData map[string]bool
}

// New returns a new BrewC instance.
func New(githubToken string, threads int) *BrewC {
	return &BrewC{
		threads:     threads,
		archOSName:  models.GetArchAndOSName(),
		githubToken: githubToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		brew:       brew.New(),
		downloader: downloader.New(),
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

	err = b.DownloadFormulas(list)

	if err != nil {
		return err
	}

	return nil
}

// GetAllFormulas returns all of the formulas.
func (b *BrewC) GetAllFormulas(name string) ([]*formula.Formula, error) {
}

// getAllFormulasRecursive returns all of the formulas recursively.
func (b *BrewC) getAllFormulasRecursive(name string, list []*formula.Formula, data map[string]int, nestedCount int) error {

	var wg sync.WaitGroup
	var conn = make(chan int, b.threads)

	var mainFormula, err = b.getFormulaJSON(name)

	if err != nil {
		return err
	}

	list = append(list, mainFormula)

	for _, dep := range mainFormula.Dependencies {
		wg.Add(1)

		go func(dep string) {
			defer wg.Done()
			conn <- 1

			f, err := b.getFormulaJSON(dep)

			if err != nil {
				fmt.Printf("error getting formula %s: %v", dep, err)
				return
			}

			fmt.Printf("discovered: %s\n", col.Info(dep))
			list = append(list, f)

			<-conn
		}(dep)
	}

	wg.Wait()

	return list, nil
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

// DownloadFormulas downloads all of the given formulas to brew's cache folder
func (b *BrewC) DownloadFormulas(list []*formula.Formula) error {
	var total = len(list)

	fmt.Printf("total formulas to download: %s\n", col.Green(total))

	for i, f := range list {
		fmt.Print(col.Green(f.Name))

		if i+1 != total {
			fmt.Print(", ")
		} else {
			fmt.Println()
		}
	}
	return nil
}

// doGET makes a GET request to the given URL.
func (b *BrewC) doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Authorization", "Bearer "+b.githubToken)

	if err != nil {
		return nil, err
	}

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
