package main

import "fmt"

// Version is the current version of llcm.
const Version = "0.0.4"

// Revision is the git revision of llcm.
var Revision = ""

// version returns the version of llcm.
func version() string {
	if Revision == "" {
		return Version
	}
	return fmt.Sprintf("%s (revision: %s)", Version, Revision)
}
