package brew

import "github.com/hamza72x/brewc/pkg/util"

type Brew struct {
	bin string
}

func New() *Brew {
	return &Brew{
		bin: getBrewBinary(),
	}
}

func getBrewBinary() string {
	var bin string

	var paths = []string{
		"/usr/local/bin/brew",
		"/opt/homebrew/bin/brew",
		"/home/linuxbrew/.linuxbrew/bin/brew",
	}

	for _, path := range paths {
		if util.DoesFileExist(path) {
			bin = path
			break
		}
	}

	if bin == "" {
		panic("brew binary not found")
	}

	return bin
}
