package version

import "fmt"

// Schema defines the version of Pactus software.
// It follows the semantic versioning 2.0.0 spec (http://semver.org/).
type Schema struct {
	Major uint   // Major version number
	Minor uint   // Minor version number
	Patch uint   // Patch version number
	Meta  string // Metadata for version (e.g., "beta", "rc1")
}

var Version = Schema{
	Major: 1,
	Minor: 0,
	Patch: 0,
	Meta:  "",
}

func (v Schema) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Meta != "" {
		version = fmt.Sprintf("%s-%s", version, v.Meta)
	}

	return version
}
