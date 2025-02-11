package llcm

import (
	"testing"
	"time"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	var (
		originalMaxPieChartItems = setMaxPieChartItems(3)
		originalMaxBarChartItems = setMaxBarChartItems(3)
		localTimeZone            = setTimeZone(time.UTC)
	)
	nowFunc = func() time.Time {
		return mustTime("2025-04-01T00:00:00Z")
	}
	defer func() {
		setMaxPieChartItems(originalMaxPieChartItems)
		setMaxBarChartItems(originalMaxBarChartItems)
		setTimeZone(localTimeZone)
		nowFunc = time.Now
	}()
	m.Run()
}

// setMaxPieChartItems is helper function to set MaxPieChartItems and return the original value.
func setMaxPieChartItems(n int) (original int) {
	original = MaxPieChartItems
	MaxPieChartItems = n
	return original
}

// setMaxBarChartItems is helper function to set MaxBarChartItems and return the original value.
func setMaxBarChartItems(n int) (original int) {
	original = MaxBarChartItems
	MaxBarChartItems = n
	return original
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
