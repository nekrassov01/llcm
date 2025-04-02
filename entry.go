package llcm

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

var (
	retentionInDaysLabel = "retentionInDays"
	storedBytesLabel     = "storedBytes"
	desiredStateLabel    = "desiredState"
	reducibleBytesLabel  = "reducibleBytes"
	remainingBytesLabel  = "remainingBytes"
)

var (
	_ Entry = (*ListEntry)(nil)
	_ Entry = (*PreviewEntry)(nil)
)

// Entry is an interface for log group entry.
type Entry interface {
	Name() string              // Name returns the name of the entry.
	Bytes() int64              // Bytes returns the stored bytes of the entry.
	DataSet() map[string]int64 // DataSet returns map for plotting the chart.
	toInput() []any            // toInput returns the input of the entry for rendering.
	toTSV() []string           // toTSV returns the TSV of the entry for rendering.
}

// entry represents the base entry for log group.
type entry struct {
	LogGroupName    string              // The name of the log group.
	Region          string              // The region that the log group belongs to.
	Source          string              // The information of linked source account.
	Class           types.LogGroupClass // The class of the log group.
	CreatedAt       time.Time           // The time when the log group was created.
	ElapsedDays     int64               // The number of days elapsed since the log group was created.
	RetentionInDays int64               // The retention days of the log group.
	StoredBytes     int64               // The stored bytes of the log group.
	name            *string             // The native type of LogGroupName.
}

// Name returns the name of the entry.
func (e *entry) Name() string {
	return e.LogGroupName
}

// Bytes returns the stored bytes of the entry.
func (e *entry) Bytes() int64 {
	return e.StoredBytes
}

// ListEntry represents an entry to list log group.
type ListEntry struct {
	*entry
}

// DataSet returns map for plotting the chart.
func (e *ListEntry) DataSet() map[string]int64 {
	return map[string]int64{
		retentionInDaysLabel: e.RetentionInDays,
		storedBytesLabel:     e.StoredBytes,
	}
}

// toInput returns the input of the list entry for rendering.
func (e *ListEntry) toInput() []any {
	return []any{
		e.LogGroupName,
		e.Region,
		e.Source,
		e.Class,
		e.CreatedAt.Format(time.RFC3339),
		e.ElapsedDays,
		e.RetentionInDays,
		e.StoredBytes,
	}
}

// toTSV returns the tab-separated values of the list entry for rendering.
func (e *ListEntry) toTSV() []string {
	return []string{
		e.LogGroupName,
		e.Region,
		e.Source,
		string(e.Class),
		e.CreatedAt.Format(time.RFC3339),
		strconv.FormatInt(e.ElapsedDays, 10),
		strconv.FormatInt(e.RetentionInDays, 10),
		strconv.FormatInt(e.StoredBytes, 10),
	}
}

// PreviewEntry is an extended representation of entry with the desired state and its simulated results.
type PreviewEntry struct {
	*entry
	BytesPerDay     int64 // The bytes per day of the log group.
	DesiredState    int64 // The desired state of the log group.
	ReductionInDays int64 // The number of days to be reduced after the action.
	ReducibleBytes  int64 // The number of bytes that can be reduced after the action.
	RemainingBytes  int64 // The number of bytes that remain after the action.
}

// DataSet returns map for plotting the chart.
func (e *PreviewEntry) DataSet() map[string]int64 {
	return map[string]int64{
		retentionInDaysLabel: e.RetentionInDays,
		storedBytesLabel:     e.StoredBytes,
		desiredStateLabel:    e.DesiredState,
		reducibleBytesLabel:  e.ReducibleBytes,
		remainingBytesLabel:  e.RemainingBytes,
	}
}

// toInput returns the input of the desired entry for rendering.
func (e *PreviewEntry) toInput() []any {
	return []any{
		e.LogGroupName,
		e.Region,
		e.Source,
		e.Class,
		e.CreatedAt.Format(time.RFC3339),
		e.ElapsedDays,
		e.RetentionInDays,
		e.StoredBytes,
		e.BytesPerDay,
		e.DesiredState,
		e.ReductionInDays,
		e.ReducibleBytes,
		e.RemainingBytes,
	}
}

// toTSV returns the tab-separated values of the desired entry for rendering.
func (e *PreviewEntry) toTSV() []string {
	return []string{
		e.LogGroupName,
		e.Region,
		e.Source,
		string(e.Class),
		e.CreatedAt.Format(time.RFC3339),
		strconv.FormatInt(e.ElapsedDays, 10),
		strconv.FormatInt(e.RetentionInDays, 10),
		strconv.FormatInt(e.StoredBytes, 10),
		strconv.FormatInt(e.BytesPerDay, 10),
		strconv.FormatInt(e.DesiredState, 10),
		strconv.FormatInt(e.ReductionInDays, 10),
		strconv.FormatInt(e.ReducibleBytes, 10),
		strconv.FormatInt(e.RemainingBytes, 10),
	}
}

// simulate calculates the simulated results for the log group.
func (e *PreviewEntry) simulate(desired DesiredState) {
	e.setDesiredState(desired)
	e.setBytesPerDay()
	e.setReductionInDays()
	e.setReducibleBytes()
	e.setRemainingBytes()
}

// setDesiredState sets the desired state for the log group.
func (e *PreviewEntry) setDesiredState(desired DesiredState) {
	e.DesiredState = int64(desired)
}

// setBytesPerDay sets the bytes per day for the log group.
func (e *PreviewEntry) setBytesPerDay() {
	if e.StoredBytes <= 0 {
		e.BytesPerDay = 0
		return
	}
	if e.ElapsedDays <= 0 {
		e.BytesPerDay = e.StoredBytes
		return
	}
	retentionInDays := e.RetentionInDays
	if retentionInDays >= e.ElapsedDays {
		retentionInDays = e.ElapsedDays
	}
	if retentionInDays <= int64(DesiredStateZero) {
		e.BytesPerDay = e.StoredBytes
		return
	}
	e.BytesPerDay = e.StoredBytes / retentionInDays
	// The minimum bytes per day is 1 when the stored bytes is greater than 0.
	if e.BytesPerDay <= 0 {
		e.BytesPerDay = 1
	}
}

// setReductionInDays sets the expected reduction in days after action.
func (e *PreviewEntry) setReductionInDays() {
	if e.StoredBytes <= 0 || e.BytesPerDay <= 0 {
		e.ReductionInDays = 0
		return
	}
	if e.DesiredState == int64(DesiredStateInfinite) {
		e.ReductionInDays = 0
		return
	}
	if e.DesiredState == int64(DesiredStateZero) {
		if e.RetentionInDays > int64(DesiredStateZero) && e.RetentionInDays < int64(DesiredStateInfinite) {
			e.ReductionInDays = e.RetentionInDays
		} else {
			// edge case: the retention days is less than or equal to 0.
			if e.ElapsedDays <= 0 {
				e.ReductionInDays = 1
				return
			}
			e.ReductionInDays = e.ElapsedDays
		}
		return
	}
	retentionInDays := e.RetentionInDays
	if retentionInDays >= e.ElapsedDays {
		retentionInDays = e.ElapsedDays
	}
	if retentionInDays > e.DesiredState {
		e.ReductionInDays = retentionInDays - e.DesiredState
		return
	}
	e.ReductionInDays = 0
}

// setReducibleBytes sets the expected reducible bytes after action.
func (e *PreviewEntry) setReducibleBytes() {
	if e.StoredBytes <= 0 || e.BytesPerDay <= 0 || e.ReductionInDays <= 0 || e.DesiredState == int64(DesiredStateInfinite) {
		e.ReducibleBytes = 0
		return
	}
	if e.DesiredState == int64(DesiredStateZero) {
		e.ReducibleBytes = e.StoredBytes
		return
	}
	e.ReducibleBytes = e.BytesPerDay * e.ReductionInDays
}

// setRemainingBytes sets the expected remaining bytes after action.
func (e *PreviewEntry) setRemainingBytes() {
	if e.ReducibleBytes > e.StoredBytes {
		e.RemainingBytes = 0
		return
	}
	e.RemainingBytes = e.StoredBytes - e.ReducibleBytes
}
