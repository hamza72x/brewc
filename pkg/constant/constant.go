package constant

import "os"

const DIR_CELLAR = "/usr/local/Cellar"

// DirDownloadsCache returns the path to the downloads cache directory of brew.
// e.g. /Users/hamza/Library/Caches/Homebrew/downloads
func DirDownloadsCache() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir + "/Library/Caches/Homebrew/downloads"
}
