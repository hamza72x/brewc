package models

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hamza72x/brewc/pkg/util"
	col "github.com/hamza72x/go-color"
)

// ArchAndCodeName represents the architecture and os version.
// used in brew files like: arm64_ventura, arm64_monterey, arm64_big_sur, ventura, monterey, big_sur, x86_64_linux
type ArchAndCodeName struct {
	Architecture string
	CodeName     string
}

const (
	Arm64   = "arm64"
	Aarch64 = "aarch64"
	X86_64  = "x86_64"
)

const (
	Ventura  = "ventura"
	Monterey = "monterey"
	BigSur   = "big_sur"
	Linux    = "linux"
)

var (
	macOsVersionToCodeName = map[string]string{
		"13": Ventura,
		"12": Monterey,
		"11": BigSur,
	}
)

func GetArchAndOSName() *ArchAndCodeName {
	data := &ArchAndCodeName{
		Architecture: getArchName(),
		CodeName:     getOSCodeName(),
	}

	if data.Architecture == Arm64 && data.CodeName == Linux {
		panic("arm64 linux is not supported yet")
	}

	fmt.Printf("%s: %s\n", col.Green("Platform"), data.Name())

	return data
}

// Name returns the name of the arch and os.
// example: arm64_ventura, arm64_monterey, arm64_big_sur, x86_64_linux
func (a *ArchAndCodeName) Name() string {
	full := a.Architecture + "_" + a.CodeName

	if runtime.GOOS == "darwin" {
		if a.Architecture == Arm64 {
			return full
		}
		return a.CodeName
	}

	return full
}

// getArchName returns the arch name.
// example: arm64, x86_64
func getArchName() string {

	archs := []string{"arm64", "aarch64", "x86_64"}

	arch := strings.ToLower(util.ExecMustWithTrim("uname", "-m"))

	// if arch is aarch64, change it to arm64
	if arch == "aarch64" {
		arch = "arm64"
	}

	if !util.StrContains(archs, arch) {
		panic("unknown arch: " + arch)
	}

	return arch
}

// getOSCodeName returns the os name.
// example: ventura, monterey, big_sur, linux
func getOSCodeName() string {

	versions := []string{"ventura", "monterey", "big_sur", "linux"}

	os := strings.ToLower(util.ExecMustWithTrim("uname", "-s"))

	if os == "darwin" {
		versionFull := util.ExecMustWithTrim("sw_vers", "-productVersion")
		version := strings.Split(versionFull, ".")

		if len(version) < 2 {
			panic("unknown os version: " + versionFull)
		}

		os = macOsVersionToCodeName[version[0]]
	}

	if !util.StrContains(versions, os) {
		panic("unknown os version: " + os)
	}

	return os
}
