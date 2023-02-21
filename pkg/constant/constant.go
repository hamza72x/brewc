package constant

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hamza72x/brewc/pkg/util"
	col "github.com/hamza72x/go-color"
)

var GreenArrow = col.Green("<==>")
var BlueArrow = col.Info("<==>")
var RedArrow = col.Red("<==>")

type Constant struct {
	DirCellar    string
	DirCaches    string
	DirDownloads string
}

var instance *Constant

func Initialize(arch string) {
	dirHome, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	dirCellar := "/usr/local/Cellar"

	if arch == "arm64" && runtime.GOOS == "darwin" {
		dirCellar = "/opt/homebrew/Cellar"
	}

	instance = &Constant{
		DirCellar:    dirCellar,
		DirCaches:    dirHome + "/Library/Caches/Homebrew",
		DirDownloads: dirHome + "/Library/Caches/Homebrew/downloads",
	}

	// create dirs
	var dirs = []string{
		instance.DirCaches,
		instance.DirDownloads,
	}

	for _, dir := range dirs {
		if util.CreateDirIfNotExists(dir) != nil {
			panic(fmt.Sprintf("Failed to create dir: %s", dir))
		}
	}
}

func Get() *Constant {
	if instance == nil {
		panic("Constant not initialized")
	}
	return instance
}
