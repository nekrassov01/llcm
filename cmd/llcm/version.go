package main

import "fmt"

// Version is the current version of llcm.
const Version = "0.0.6"

// Commit is the git revision of llcm.
var Commit = ""

// version returns the version and revision of llcm.
func version() string {
	if Commit == "" {
		return Version
	}
	return fmt.Sprintf("%s (revision: %s)", Version, Commit)
}
