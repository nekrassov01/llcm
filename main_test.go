package llcm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
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
		if err := removeChartFiles(); err != nil {
			panic(err)
		}
	}()
	m.Run()
}

func removeChartFiles() error {
	files, err := filepath.Glob("llcm*.html")
	if err != nil {
		return err
	}
	for _, f := range files {
		if strings.HasPrefix(filepath.Base(f), "llcm") && strings.HasSuffix(f, ".html") {
			if err := os.Remove(f); err != nil {
				fmt.Printf("failed to remove file %s: %v\n", f, err)
			}
		}
	}
	return nil
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

// errData is a test data for ListEntryData of error case.
var errData = ListEntryData{
	header: previewEntryDataHeader,
	entries: []*ListEntry{
		{
			entry: &entry{},
		},
	},
}

// listEntryData is a test data for ListEntryData.
var listEntryData = ListEntryData{
	header: listEntryDataHeader,
	entries: []*ListEntry{
		{
			entry: &entry{
				LogGroupName:    "group0",
				Region:          "ap-northeast-1",
				Class:           types.LogGroupClassStandard,
				CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
				ElapsedDays:     90,
				RetentionInDays: 30,
				StoredBytes:     1024,
				name:            aws.String("group0"),
			},
		},
		{
			entry: &entry{
				LogGroupName:    "group1",
				Region:          "ap-northeast-2",
				Class:           types.LogGroupClassInfrequentAccess,
				CreatedAt:       mustTime("2024-04-01T00:00:00Z"),
				ElapsedDays:     365,
				RetentionInDays: 30,
				StoredBytes:     2048,
				name:            aws.String("group1"),
			},
		},
	},
}

// previewEntryData is a test data for PreviewEntryData.
var previewEntryData = PreviewEntryData{
	header: previewEntryDataHeader,
	entries: []*PreviewEntry{
		{
			BytesPerDay:     0,
			DesiredState:    0,
			ReductionInDays: 0,
			ReducibleBytes:  0,
			RemainingBytes:  0,
			entry: &entry{
				LogGroupName:    "group0",
				Region:          "ap-northeast-1",
				Class:           types.LogGroupClassStandard,
				CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
				ElapsedDays:     90,
				RetentionInDays: 30,
				StoredBytes:     1024,
				name:            aws.String("group0"),
			},
		},
		{
			BytesPerDay:     100,
			DesiredState:    100,
			ReductionInDays: 100,
			ReducibleBytes:  100,
			RemainingBytes:  100,
			entry: &entry{
				LogGroupName:    "group1",
				Region:          "ap-northeast-2",
				Class:           types.LogGroupClassInfrequentAccess,
				CreatedAt:       mustTime("2024-04-01T00:00:00Z"),
				ElapsedDays:     365,
				RetentionInDays: 30,
				StoredBytes:     2048,
				name:            aws.String("group1"),
			},
		},
	},
}
