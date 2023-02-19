package manifest

// https://ghcr.io/v2/homebrew/core/libraw/manifests/0.21.1
// Request Headers:
//
//		Accept: application/vnd.oci.image.index.v1+json
//	 	Authorization: Bearer <token>
type Manifest struct {
	SchemaVersion int64               `json:"schemaVersion"`
	Manifests     []ManifestElement   `json:"manifests"`
	Annotations   ManifestAnnotations `json:"annotations"`
}

type ManifestAnnotations struct {
	COMGithubPackageType                string `json:"com.github.package.type"`
	OrgOpencontainersImageCreated       string `json:"org.opencontainers.image.created"`
	OrgOpencontainersImageDescription   string `json:"org.opencontainers.image.description"`
	OrgOpencontainersImageDocumentation string `json:"org.opencontainers.image.documentation"`
	OrgOpencontainersImageLicense       string `json:"org.opencontainers.image.license"`
	OrgOpencontainersImageRefName       string `json:"org.opencontainers.image.ref.name"`
	OrgOpencontainersImageRevision      string `json:"org.opencontainers.image.revision"`
	OrgOpencontainersImageSource        string `json:"org.opencontainers.image.source"`
	OrgOpencontainersImageTitle         string `json:"org.opencontainers.image.title"`
	OrgOpencontainersImageURL           string `json:"org.opencontainers.image.url"`
	OrgOpencontainersImageVendor        string `json:"org.opencontainers.image.vendor"`
	OrgOpencontainersImageVersion       string `json:"org.opencontainers.image.version"`
}

type ManifestElement struct {
	MediaType   string                   `json:"mediaType"`
	Digest      string                   `json:"digest"`
	Size        int64                    `json:"size"`
	Platform    Platform                 `json:"platform"`
	Annotations ManifestAnnotationsClass `json:"annotations"`
}

type ManifestAnnotationsClass struct {
	OrgOpencontainersImageRefName string  `json:"org.opencontainers.image.ref.name"`
	ShBrewBottleDigest            string  `json:"sh.brew.bottle.digest"`
	ShBrewTab                     string  `json:"sh.brew.tab"`
	ShBrewBottleCPUVariant        *string `json:"sh.brew.bottle.cpu.variant,omitempty"`
	ShBrewBottleGlibcVersion      *string `json:"sh.brew.bottle.glibc.version,omitempty"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	OSVersion    string `json:"os.version"`
}
