package llcm

var (
	_ EntryData[*ListEntry]    = (*ListEntryData)(nil)
	_ EntryData[*PreviewEntry] = (*PreviewEntryData)(nil)
)

var (
	// TotalStoredBytesLabel is the label of the total stored bytes.
	TotalStoredBytesLabel = "storedBytes"

	// TotalReducibleBytesLabel is the label of the total reducible bytes.
	TotalReducibleBytesLabel = "reducibleBytes"

	// TotalRemainingBytesLabel is the label of the total remaining bytes.
	TotalRemainingBytesLabel = "remainingBytes"
)

var (
	// listEntryDataHeader is the header of ListEntryData.
	listEntryDataHeader = []string{
		"Name",
		"Region",
		"Source",
		"Class",
		"CreatedAt",
		"ElapsedDays",
		"RetentionInDays",
		"StoredBytes",
	}

	// previewEntryDataHeader is the header of PreviewEntryData.
	previewEntryDataHeader = []string{
		"Name",
		"Region",
		"Source",
		"Class",
		"CreatedAt",
		"ElapsedDays",
		"RetentionInDays",
		"StoredBytes",
		"BytesPerDay",
		"DesiredState",
		"ReductionInDays",
		"ReducibleBytes",
		"RemainingBytes",
	}
)

// EntryData represents the collection of entries.
type EntryData[T Entry] interface {
	Header() []string
	Entries() []T
	Total() map[string]int64
}

// ListEntryData represents the collection of ListEntry.
type ListEntryData struct {
	TotalStoredBytes int64 // TotalStoredBytes represents the total stored bytes of the log groups.

	header  []string
	entries []*ListEntry
}

// Header returns the header of the ListEntryData.
func (d *ListEntryData) Header() []string {
	return d.header
}

// Entries returns the entries of the ListEntryData.
func (d *ListEntryData) Entries() []*ListEntry {
	return d.entries
}

// Total returns the total of the ListEntryData.
func (d *ListEntryData) Total() map[string]int64 {
	return map[string]int64{
		TotalStoredBytesLabel: d.TotalStoredBytes,
	}
}

// PreviewEntryData represents the collection of PreviewEntry.
type PreviewEntryData struct {
	TotalStoredBytes    int64 // TotalStoredBytes represents the total stored bytes of the log groups.
	TotalReducibleBytes int64 // TotalReducibleBytes represents the total reducible bytes of the log groups.
	TotalRemainingBytes int64 // TotalRemainingBytes represents the total remaining bytes of the log groups.

	header  []string
	entries []*PreviewEntry
}

// Header returns the header of the PreviewEntryData.
func (d *PreviewEntryData) Header() []string {
	return d.header
}

// Entries returns the entries of the PreviewEntryData.
func (d *PreviewEntryData) Entries() []*PreviewEntry {
	return d.entries
}

// Total returns the total of the PreviewEntryData.
func (d *PreviewEntryData) Total() map[string]int64 {
	return map[string]int64{
		TotalStoredBytesLabel:    d.TotalStoredBytes,
		TotalReducibleBytesLabel: d.TotalReducibleBytes,
		TotalRemainingBytesLabel: d.TotalRemainingBytes,
	}
}
