package models

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hamza72x/brewc/pkg/util"
)

// ArchAndOS represents the architecture and os version.
// used in brew files like: arm64_ventura, arm64_monterey, arm64_big_sur, ventura, monterey, big_sur, x86_64_linux
type ArchAndOS struct {
	Architecture Architecture
	OSVersion    OSVersion
}

type Architecture string
type OSVersion string

const (
	Arm64   Architecture = "arm64"
	Aarch64 Architecture = "aarch64"
	X86_64  Architecture = "x86_64"
)

const (
	Ventura  OSVersion = "ventura"
	Monterey OSVersion = "monterey"
	BigSur   OSVersion = "big_sur"
	Linux    OSVersion = "linux"
)

func GetArchAndOSName() *ArchAndOS {
	data := &ArchAndOS{
		Architecture: getArchName(),
		OSVersion:    getOSName(),
	}

	if data.Architecture == Arm64 && data.OSVersion == Linux {
		// panic("arm64 linux is not supported yet")
	}

	fmt.Printf("Detected arch and os: %s\n", data.Name())

	return data
}

// Name returns the name of the arch and os.
// example: arm64_ventura, arm64_monterey, arm64_big_sur, x86_64_linux
func (a *ArchAndOS) Name() string {
	full := string(a.Architecture) + "_" + string(a.OSVersion)

	if runtime.GOOS == "darwin" {
		if a.Architecture == Arm64 {
			return full
		}
		return string(a.OSVersion)
	}

	return full
}

// getArchName returns the arch name.
// example: arm64, x86_64
func getArchName() Architecture {

	archs := []string{"arm64", "aarch64", "x86_64"}

	res, err := util.Exec("uname", "-m")

	if err != nil {
		panic(err)
	}

	arch := strings.ToLower(strings.TrimSpace(res))

	// if arch is aarch64, change it to arm64
	if arch == "aarch64" {
		arch = "arm64"
	}

	if !util.StrContains(archs, arch) {
		panic("unknown arch: " + arch)
	}

	return Architecture(arch)
}

// getOSName returns the os name.
// example: ventura, monterey, big_sur, linux
func getOSName() OSVersion {

	versions := []string{"ventura", "monterey", "big_sur", "linux"}

	res, err := util.Exec("uname", "-s")

	if err != nil {
		panic(err)
	}

	os := strings.ToLower(strings.TrimSpace(res))

	if !util.StrContains(versions, os) {
		panic("unknown os version: " + os)
	}

	return OSVersion(os)
}
