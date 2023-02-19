package downloader

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hamza72x/brewc/pkg/models"
	"github.com/hamza72x/brewc/pkg/models/formula"
	"github.com/hamza72x/brewc/pkg/models/manifest"
)

type Downloader struct {
	githubToken     string
	archAndCodeName *models.ArchAndCodeName
	client          *http.Client
}

func New(archAndCodeName *models.ArchAndCodeName, githubToken string) *Downloader {
	return &Downloader{
		archAndCodeName: archAndCodeName,
		githubToken:     githubToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Download downloads the file from the given url
// and saves it to the given path.
func (d *Downloader) Download(f *formula.Formula) error {
	m, err := d.getManifest(f)

	if err != nil {
		return err
	}

	if !f.HasBottleDownloadCache(d.archAndCodeName.CodeName) {

	}

	return nil
}

func (d *Downloader) downloadFormula(f *formula.Formula) error {
	// req, err := http.NewRequest(http.MethodGet, url, nil)

	// if err != nil {
	// 	return err
	// }

	// resp, err := d.client.Do(req)

	// if err != nil {
	// 	return err
	// }

	// defer resp.Body.Close()

	// return util.WriteFile(toPath, resp.Body)
	return nil
}

func (d *Downloader) getManifest(f *formula.Formula) (*manifest.Manifest, error) {
	if f.HasManifestDownloadCache() {
		// var m manifest.Manifest
	}

	return d.downloadManifest(f)
}

func (d *Downloader) downloadManifest(f *formula.Formula) (*manifest.Manifest, error) {
	req, err := http.NewRequest(http.MethodGet, f.GetManifestUrl(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")
	req.Header.Set("Authorization", "Bearer "+d.githubToken)

	resp, err := d.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var m manifest.Manifest

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
