package formula

// Formula represents a formula.
// GET https://formulae.brew.sh/api/formula/${FORMULA}.json
// Example: curl -sL https://formulae.brew.sh/api/formula/ffmpeg.json | jq
type Formula struct {
	Name                    string         `json:"name"`
	FullName                string         `json:"full_name"`
	Tap                     string         `json:"tap"`
	Oldname                 *string        `json:"oldname"`
	Aliases                 *[]string      `json:"aliases"`
	VersionedFormulae       *[]string      `json:"versioned_formulae"`
	Desc                    string         `json:"desc"`
	License                 string         `json:"license"`
	Homepage                string         `json:"homepage"`
	Versions                Versions       `json:"versions"`
	Urls                    Urls           `json:"urls"`
	Revision                int64          `json:"revision"`
	VersionScheme           int64          `json:"version_scheme"`
	Bottle                  Bottle         `json:"bottle"`
	KegOnly                 bool           `json:"keg_only"`
	KegOnlyReason           *string        `json:"keg_only_reason"`
	Options                 *[]string      `json:"options"`
	BuildDependencies       []string       `json:"build_dependencies"`
	Dependencies            []string       `json:"dependencies"`
	TestDependencies        *[]string      `json:"test_dependencies"`
	RecommendedDependencies *[]string      `json:"recommended_dependencies"`
	OptionalDependencies    *[]string      `json:"optional_dependencies"`
	UsesFromMacos           *[]string      `json:"uses_from_macos"`
	Requirements            *[]string      `json:"requirements"`
	ConflictsWith           *[]string      `json:"conflicts_with"`
	Caveats                 *string        `json:"caveats"`
	Installed               []Installed    `json:"installed"`
	LinkedKeg               string         `json:"linked_keg"`
	Pinned                  bool           `json:"pinned"`
	Outdated                bool           `json:"outdated"`
	Deprecated              bool           `json:"deprecated"`
	DeprecationDate         *string        `json:"deprecation_date"`
	DeprecationReason       *string        `json:"deprecation_reason"`
	Disabled                bool           `json:"disabled"`
	DisableDate             *string        `json:"disable_date"`
	DisableReason           *string        `json:"disable_reason"`
	TapGitHead              string         `json:"tap_git_head"`
	Variations              Variations     `json:"variations"`
	Analytics               Analytics      `json:"analytics"`
	AnalyticsLinux          AnalyticsLinux `json:"analytics-linux"`
	GeneratedDate           string         `json:"generated_date"`
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
	Arm64Ventura  Arm64BigSur `json:"arm64_ventura"`
	Arm64Monterey Arm64BigSur `json:"arm64_monterey"`
	Arm64BigSur   Arm64BigSur `json:"arm64_big_sur"`
	Ventura       Arm64BigSur `json:"ventura"`
	Monterey      Arm64BigSur `json:"monterey"`
	BigSur        Arm64BigSur `json:"big_sur"`
	X8664_Linux   Arm64BigSur `json:"x86_64_linux"`
}

type Arm64BigSur struct {
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
