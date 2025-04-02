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
	Chart() error
}

// ListEntryData represents the collection of ListEntry.
type ListEntryData struct {
	TotalStoredBytes int64 // The total stored bytes of the log groups.

	header  []string
	entries []*ListEntry
}

// Header returns the header of the ListEntryData.
func (d *ListEntryData) Header() []string {
	return d.header
}

// Entries returns the entries of the ListEntryData.
func (d *ListEntryData) Entries() []*ListEntry {
	if len(d.entries) == 0 {
		return nil
	}
	return d.entries
}

// Total returns the total of the ListEntryData.
func (d *ListEntryData) Total() map[string]int64 {
	return map[string]int64{
		TotalStoredBytesLabel: d.TotalStoredBytes,
	}
}

// Chart generates a pie chart for the ListEntryData.
func (d *ListEntryData) Chart() error {
	if len(d.entries) == 0 {
		return nil
	}
	items := getPieItems(d.entries)
	chart := newPieChart(items)
	if chart == nil {
		return nil
	}
	return render(chart)
}

// PreviewEntryData represents the collection of PreviewEntry.
type PreviewEntryData struct {
	TotalStoredBytes    int64 // The total stored bytes of the log groups.
	TotalReducibleBytes int64 // The total reducible bytes of the log groups.
	TotalRemainingBytes int64 // The total remaining bytes of the log groups.

	header  []string
	entries []*PreviewEntry
}

// Header returns the header of the PreviewEntryData.
func (d *PreviewEntryData) Header() []string {
	return d.header
}

// Entries returns the entries of the PreviewEntryData.
func (d *PreviewEntryData) Entries() []*PreviewEntry {
	if len(d.entries) == 0 {
		return nil
	}
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

// Chart generates a bar chart for the PreviewEntryData.
func (d *PreviewEntryData) Chart() error {
	if len(d.entries) == 0 {
		return nil
	}
	subtitle := getBarSubtitle(d.entries)
	lnames, rmbytes, rdbytes := getBarItems(d.entries)
	chart := newBarChart(subtitle, lnames, rmbytes, rdbytes)
	if chart == nil {
		return nil
	}
	return render(chart)
}
