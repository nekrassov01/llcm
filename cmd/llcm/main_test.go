package main

import (
	"testing"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	revision := Revision
	defer func() {
		Revision = revision
	}()
	m.Run()
}
