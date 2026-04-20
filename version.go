package llcm

import "fmt"

// version is the current version.
const version = "0.0.30"

// revision is the git revision.
var revision = ""

// Version returns the version and revision.
func Version() string {
	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (revision: %s)", version, revision)
}
