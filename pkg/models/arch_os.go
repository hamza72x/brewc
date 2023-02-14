package models

import (
	"strings"

	"github.com/hamza72x/brewc/pkg/util"
)

// ArchOSName, example: arm64_ventura, arm64_monterey, arm64_big_sur, ventura, monterey, big_sur, x86_64_linux
type ArchOSName struct {
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
	OSVersionVentura  OSVersion = "ventura"
	OSVersionMonterey OSVersion = "monterey"
	OSVersionBigSur   OSVersion = "big_sur"
	OSVersionLinux    OSVersion = "linux"
)

func GetArchAndOSName() *ArchOSName {
	data := &ArchOSName{
		Architecture: getArchName(),
		OSVersion:    getOSName(),
	}

	if data.Architecture == Arm64 && data.OSVersion == OSVersionLinux {
		panic("arm64 linux is not supported yet")
	}

	return data
}

// getArchName returns the arch name.
// example: arm64, x86_64
func getArchName() Architecture {
	res, err := util.Exec("uname", "-m")

	if err != nil {
		panic(err)
	}

	arch := Architecture(strings.ToLower(strings.TrimSpace(res)))

	if arch == Aarch64 {
		arch = Arm64
	}

	// assert
	switch arch {
	case Arm64, X86_64:
		break
	default:
		panic("unknown arch: " + arch)
	}

	return Architecture(arch)
}

// getOSName returns the os name.
// example: ventura, monterey, big_sur, linux
func getOSName() OSVersion {
	res, err := util.Exec("uname", "-s")

	if err != nil {
		panic(err)
	}

	os := OSVersion(strings.ToLower(strings.TrimSpace(res)))

	// assert
	switch os {
	case OSVersionVentura, OSVersionMonterey, OSVersionBigSur, OSVersionLinux:
		break
	default:
		panic("unknown os: " + os)
	}

	return os
}
