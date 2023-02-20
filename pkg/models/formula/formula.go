package formula

import (
	"fmt"

	"github.com/hamza72x/brewc/pkg/constant"
	"github.com/hamza72x/brewc/pkg/util"
)

// Formula represents a formula.
// GET https://formulae.brew.sh/api/formula/${FORMULA}.json
// Example: curl -sL https://formulae.brew.sh/api/formula/ffmpeg.json | jq
type Formula struct {
	Name                    string    `json:"name"`
	FullName                string    `json:"full_name"`
	Tap                     string    `json:"tap"`
	Desc                    string    `json:"desc"`
	Versions                Versions  `json:"versions"`
	Urls                    Urls      `json:"urls"`
	Revision                int64     `json:"revision"`
	VersionScheme           int64     `json:"version_scheme"`
	Bottle                  Bottle    `json:"bottle"`
	BuildDependencies       []string  `json:"build_dependencies"`
	Dependencies            []string  `json:"dependencies"`
	TestDependencies        *[]string `json:"test_dependencies"`
	RecommendedDependencies *[]string `json:"recommended_dependencies"`
	OptionalDependencies    *[]string `json:"optional_dependencies"`
	Requirements            *[]string `json:"requirements"`
	ConflictsWith           *[]string `json:"conflicts_with"`
	Caveats                 *string   `json:"caveats"`
	Outdated                bool      `json:"outdated"`
	Deprecated              bool      `json:"deprecated"`
	DeprecationDate         *string   `json:"deprecation_date"`
	DeprecationReason       *string   `json:"deprecation_reason"`
	Disabled                bool      `json:"disabled"`
	DisableDate             *string   `json:"disable_date"`
	DisableReason           *string   `json:"disable_reason"`
	GeneratedDate           string    `json:"generated_date"`
}

type Analytics struct {
	Install          map[string]Analytics30D `json:"install"`
	InstallOnRequest map[string]Analytics30D `json:"install_on_request"`
	BuildError       AnalyticsBuildError     `json:"build_error"`
}

type AnalyticsBuildError struct {
	The30D Analytics30D `json:"30d"`
}

type Analytics30D struct {
	Wget     int64 `json:"wget"`
	WgetHEAD int64 `json:"wget --HEAD"`
}

type AnalyticsLinux struct {
	Install          map[string]AnalyticsLinux30D `json:"install"`
	InstallOnRequest map[string]AnalyticsLinux30D `json:"install_on_request"`
	BuildError       AnalyticsLinuxBuildError     `json:"build_error"`
}

type AnalyticsLinuxBuildError struct {
	The30D AnalyticsLinux30D `json:"30d"`
}

type AnalyticsLinux30D struct {
	Wget int64 `json:"wget"`
}

type Bottle struct {
	Stable BottleStable `json:"stable"`
}

type BottleStable struct {
	Rebuild int64  `json:"rebuild"`
	RootURL string `json:"root_url"`
	Files   Files  `json:"files"`
}

type Files struct {
	Arm64Ventura  BottleUrlData `json:"arm64_ventura"`
	Arm64Monterey BottleUrlData `json:"arm64_monterey"`
	Arm64BigSur   BottleUrlData `json:"arm64_big_sur"`
	Ventura       BottleUrlData `json:"ventura"`
	Monterey      BottleUrlData `json:"monterey"`
	BigSur        BottleUrlData `json:"big_sur"`
	X8664_Linux   BottleUrlData `json:"x86_64_linux"`
}

type BottleUrlData struct {
	Cellar string `json:"cellar"`
	URL    string `json:"url"`
	Sha256 string `json:"sha256"`
}

type Installed struct {
	Version               string              `json:"version"`
	UsedOptions           *[]string           `json:"used_options"`
	BuiltAsBottle         bool                `json:"built_as_bottle"`
	PouredFromBottle      bool                `json:"poured_from_bottle"`
	Time                  int64               `json:"time"`
	RuntimeDependencies   []RuntimeDependency `json:"runtime_dependencies"`
	InstalledAsDependency bool                `json:"installed_as_dependency"`
	InstalledOnRequest    bool                `json:"installed_on_request"`
}

type RuntimeDependency struct {
	FullName         string `json:"full_name"`
	Version          string `json:"version"`
	DeclaredDirectly bool   `json:"declared_directly"`
}

type Urls struct {
	Stable UrlsStable `json:"stable"`
	Head   Head       `json:"head"`
}

type Head struct {
	URL    string  `json:"url"`
	Branch *string `json:"branch"`
}

type UrlsStable struct {
	URL      string  `json:"url"`
	Tag      *string `json:"tag"`
	Revision *string `json:"revision"`
	Checksum string  `json:"checksum"`
}

type UsesFromMacoClass struct {
	Gperf  *string `json:"gperf,omitempty"`
	Python *string `json:"python,omitempty"`
}

type Variations struct {
	X8664_Linux X8664_Linux `json:"x86_64_linux"`
}

type X8664_Linux struct {
	Dependencies []string `json:"dependencies"`
}

type Versions struct {
	Stable string `json:"stable"`
	Head   string `json:"head"`
	Bottle bool   `json:"bottle"`
}

type UsesFromMacoElement struct {
	String            *string
	UsesFromMacoClass *UsesFromMacoClass
}

// IsInstalled returns true if the formula is installed
// based on the folder existence of /usr/local/Cellar/{name}/{version}
func (f *Formula) IsInstalled() bool {
	return util.DoesDirExist(fmt.Sprintf("%s/%s/%s", constant.Get().DirCellar, f.Name, f.Versions.Stable))
}

// GetBottleUrl returns the bottle url of the formula
// example: https://ghcr.io/v2/homebrew/core/libraw/blobs/sha256:81a83bd632b57ca84ce11f0829942a8061c7a57d3568e6c20c54c919fa2c6111
func (f *Formula) GetBottleUrl(osCodeName string) string {
	files := f.Bottle.Stable.Files

	switch osCodeName {
	case "arm64_ventura":
		return files.Arm64Ventura.URL
	case "arm64_monterey":
		return files.Arm64Monterey.URL
	case "arm64_big_sur":
		return files.Arm64BigSur.URL
	case "ventura":
		return files.Ventura.URL
	case "monterey":
		return files.Monterey.URL
	case "big_sur":
		return files.BigSur.URL
	case "x86_64_linux":
		return files.X8664_Linux.URL
	default:
		panic(fmt.Sprintf("unknown os code name: %s", osCodeName))
	}
}

// HasBottleDownloadCache returns true if the bottle download cache exists
func (f *Formula) HasBottleDownloadCache(osCodeName string) bool {
	return util.DoesFileExist(f.GetBottleDownloadPath(osCodeName))
}

// GetBottleDownloadPath returns the cache path of the bottle
// example: $HOMEBREW_DOWNLOADS_DIR/ff7fbec7b5a2946b14760f437f4e71201b7d0bdf2d68ebdcf4d308eece3e5061--luajit--2.1.0-beta3-20230104.2.ventura.bottle.tar.gz
func (f *Formula) GetBottleDownloadPath(osCodeName string) string {
	// sha256 of url
	url := util.Sha256(f.GetBottleUrl(osCodeName))

	// example:
	// ff7fbec7b5a2946b14760f437f4e71201b7d0bdf2d68ebdcf4d308eece3e5061--luajit--2.1.0-beta3-20230104.2.ventura.bottle.tar.gz
	// sha256_of_url--name--version.os_code_name.bottle.tar.gz

	return fmt.Sprintf("%s/%s--%s--%s.%s.bottle.tar.gz", constant.Get().DirDownloads, string(url[:]), f.Name, f.Versions.Stable, osCodeName)
}

// GetBottleAliasPath returns the alias path of the bottle
// example: $HOMEBREW_DOWNLOADS_DIR/aribb24--1.0.4
func (f *Formula) GetBottleAliasPath() string {
	// example:
	// aribb24--1.0.4
	// name--version

	return fmt.Sprintf("%s/%s--%s", constant.Get().DirCaches, f.Name, f.Versions.Stable)
}

// GetManifestUrl returns the manifest url of the formula
// example: https://ghcr.io/v2/homebrew/core/libraw/manifests/0.21.1
func (f *Formula) GetManifestUrl() string {
	// example: https://ghcr.io/v2/homebrew/core/libraw/manifests/0.21.1
	return fmt.Sprintf("https://ghcr.io/v2/homebrew/core/%s/manifests/%s", f.Name, f.Versions.Stable)
}

// GetManifestDownloadPath returns the cache path of the manifest
// example: $HOMEBREW_DOWNLOADS_DIR/dce2f2976851d7b9a08cc4fb5bcc12aab7cf40bbdfec362ef68672a15fa47e55--libvmaf-2.3.1.bottle_manifest.json
func (f *Formula) GetManifestDownloadPath() string {
	url := util.Sha256(f.GetManifestUrl())

	// example:
	// dce2f2976851d7b9a08cc4fb5bcc12aab7cf40bbdfec362ef68672a15fa47e55--libvmaf-2.3.1.bottle_manifest.json
	// sha256_of_url--name--version.bottle_manifest.json

	return fmt.Sprintf("%s/%s--%s-%s.bottle_manifest.json", constant.Get().DirDownloads, url[:], f.Name, f.Versions.Stable)
}

// GetManifestAliasPath returns the alias path of the manifest
// example: $HOMEBREW_DOWNLOADS_DIR/aribb24_bottle_manifest--1.0.4
func (f *Formula) GetManifestAliasPath() string {
	// example:
	// aribb24_bottle_manifest--1.0.4
	// name_bottle_manifest--version
	return fmt.Sprintf("%s/%s_bottle_manifest--%s", constant.Get().DirCaches, f.Name, f.Versions.Stable)
}

// HasManifestDownloadCache returns true if the manifest download cache exists
func (f *Formula) HasManifestDownloadCache() bool {
	return util.DoesFileExist(f.GetManifestDownloadPath())
}
