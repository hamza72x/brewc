package downloader

import (
	"net/http"
	"time"

	"github.com/hamza72x/brewc/pkg/util"
)

type Downloader struct {
	client *http.Client
}

func New() *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Download downloads the file from the given url
// and saves it to the given path.
func (d *Downloader) Download(url string, toPath string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return util.WriteFile(toPath, resp.Body)
}
