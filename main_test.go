package llcm

import (
	"os"
	"testing"
	"time"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	local := setTimeZone(time.UTC)
	nowFunc = func() time.Time {
		return mustTime("2025-04-01T00:00:00Z")
	}
	code := m.Run()
	defer func() {
		_ = setTimeZone(local)
		nowFunc = time.Now
		os.Exit(code)
	}()
}

// setTimeZone is helper function to set time zone and return the original time zone.
func setTimeZone(z *time.Location) (original *time.Location) {
	original = time.Local
	time.Local = z
	return original
}

// mustTime is helper function to parse time string to time.Time.
// It panics if the string is not in RFC3339 format.
func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// mustUnixMilli is helper function to parse time string to Unix milliseconds.
// It panics if the string is not in RFC3339 format.
func mustUnixMilli(s string) int64 {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.UnixMilli()
}
