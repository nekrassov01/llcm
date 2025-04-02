package main

import (
	"testing"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	original := revision
	defer func() {
		revision = original
	}()
	m.Run()
}
