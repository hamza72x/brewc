package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hamza72x/brewc/pkg/models"
	"github.com/hamza72x/brewc/pkg/models/formula"
	"github.com/hamza72x/brewc/pkg/models/manifest"
	col "github.com/hamza72x/go-color"
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
			Timeout: 30 * time.Minute,
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
		return d.downloadFormula(f, m)
	}

	return nil
}

func (d *Downloader) downloadFormula(f *formula.Formula, m *manifest.Manifest) error {
	url := f.GetBottleUrl(d.archAndCodeName.CodeName)
	bottlePath := f.GetBottleDownloadPath(d.archAndCodeName.CodeName)

	fmt.Println("Downloading", f.Name, "from", col.Info(url))

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer QQ==")
	// req.Header.Set("Authorization", "Bearer "+d.githubToken)

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s, status code: %d", f.Name, resp.StatusCode)
	}

	// copy response body to file
	file, err := os.Create(bottlePath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	// link alias to caches dir
	if err := os.Symlink(bottlePath, f.GetBottleAliasPath()); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) getManifest(f *formula.Formula) (*manifest.Manifest, error) {
	if !f.HasManifestDownloadCache() {
		if err := d.downloadManifest(f); err != nil {
			return nil, err
		}
	}

	return d.getManifestFromCache(f)
}

func (d *Downloader) downloadManifest(f *formula.Formula) error {
	url := f.GetManifestUrl()
	manifestPath := f.GetManifestDownloadPath()

	fmt.Println("Downloading manifest for", f.Name, "from", col.Info(url))

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")
	req.Header.Set("Authorization", "Bearer QQ==")
	// req.Header.Set("Authorization", "Bearer "+d.githubToken)

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download manifest for %s, status code: %d", f.Name, resp.StatusCode)
	}

	// copy response body to file
	file, err := os.Create(manifestPath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	// link alias to caches dir
	if err := os.Symlink(manifestPath, f.GetManifestAliasPath()); err != nil {
		return err
	}

	return err
}

func (d *Downloader) getManifestFromCache(f *formula.Formula) (*manifest.Manifest, error) {
	var m manifest.Manifest

	var fileBytes, err = os.ReadFile(f.GetManifestDownloadPath())

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fileBytes, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
