package main

import "fmt"

// version is the current version of llcm.
const version = "0.0.28"

// revision is the git revision of llcm.
var revision = ""

// getVersion returns the version and revision of llcm.
func getVersion() string {
	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (revision: %s)", version, revision)
}
