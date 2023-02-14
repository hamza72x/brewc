package brewc

import (
	"net/http"
	"time"

	"github.com/hamza72x/brewc/pkg/brew"
	"github.com/hamza72x/brewc/pkg/downloader"
	"github.com/hamza72x/brewc/pkg/models"
	"github.com/hamza72x/brewc/pkg/models/formula"
)

// BrewC downloads all of the dependencies for a formula in concurrent goroutines.
// and then handles to `brew` command to install the formula.
type BrewC struct {
	archOSName *models.ArchOSName

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
}

// New returns a new BrewC instance.
func New(githubToken string) *BrewC {
	return &BrewC{
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
	return nil
}

func (b *BrewC) GetAllFormulas(name string) ([]formula.Formula, error) {
	var list []formula.Formula
	return list, nil
}
