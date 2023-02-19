package formula

import (
	"fmt"
	"sync"

	"github.com/hamza72x/brewc/pkg/constant"
	"github.com/hamza72x/brewc/pkg/util"
)

// FormulaList represents a list of formulas.
// it will be linked-list
type FormulaList struct {
	head *FormulaNode
	tail *FormulaNode
	// key string: formula name
	hasDataMap map[string]bool
	count      int
	// lock is used to make the list thread-safe
	lock sync.RWMutex
}

// FormulaNode represents a node in the linked-list
type FormulaNode struct {
	Formula *Formula
	Next    *FormulaNode
}

// NewFormulaList returns a new FormulaList instance.
func NewFormulaList() *FormulaList {
	return &FormulaList{
		hasDataMap: make(map[string]bool),
	}
}

func (list *FormulaList) HasFormula(f *Formula) bool {
	if _, ok := list.hasDataMap[f.Name]; ok {
		return true
	}
	return false
}

func (list *FormulaList) Count() int {
	return list.count
}

func (list *FormulaList) Add(formula *Formula) {
	list.lock.Lock()
	defer list.lock.Unlock()

	newNode := &FormulaNode{
		Formula: formula,
	}

	if list.head == nil {
		list.head = newNode
		list.tail = newNode
	} else {
		list.tail.Next = newNode
		list.tail = newNode
	}

	list.hasDataMap[formula.Name] = true
	list.count++
}

// Iterate iterate over the full linked-list
// and uses the callback function to do something
// with the data
func (list *FormulaList) Iterate(callback func(index int, formula *Formula)) {
	current := list.head
	index := 0

	for current != nil {
		callback(index, current.Formula)
		current = current.Next
		index++
	}
}

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

func (f *Formula) GetManifestUrl() string {
	// example: https://ghcr.io/v2/homebrew/core/libraw/manifests/0.21.1
	return fmt.Sprintf("https://ghcr.io/v2/homebrew/core/%s/manifests/%s", f.Name, f.Versions.Stable)
}

func (f *Formula) GetBottleCachePath(osCodeName string) string {
	// sha256 of url
	url := util.Sha256(f.GetBottleUrl(osCodeName))

	// example:
	// ff7fbec7b5a2946b14760f437f4e71201b7d0bdf2d68ebdcf4d308eece3e5061--luajit--2.1.0-beta3-20230104.2.ventura.bottle.tar.gz
	// sha256_of_url--name--version.os_code_name.bottle.tar.gz

	return fmt.Sprintf("%s/%s--%s--%s.%s.bottle.tar.gz", constant.DirDownloadsCache(), string(url[:]), f.Name, f.Versions.Stable, osCodeName)
}

func (f *Formula) GetManifestCachePath() string {
	url := util.Sha256(f.GetManifestUrl())

	// example:
	// dce2f2976851d7b9a08cc4fb5bcc12aab7cf40bbdfec362ef68672a15fa47e55--libvmaf-2.3.1.bottle_manifest.json
	// sha256_of_url--name--version.bottle_manifest.json

	return fmt.Sprintf("%s/%s--%s-%s.bottle_manifest.json", constant.DirDownloadsCache(), url[:], f.Name, f.Versions.Stable)
}

func (f *Formula) IsInstalled() bool {
	return util.DoesDirExist(fmt.Sprintf("%s/%s/%s", constant.DIR_CELLAR, f.Name, f.Versions.Stable))
}

func (f *Formula) HasDownloadCache() bool {
	// dirCahe := constant.DirDownloadsCache()

	return false
}
