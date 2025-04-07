package llcm

import (
	"sync"
)

// List lists the log group entries.
func (man *Manager) List() (*ListEntryData, error) {
	var (
		total int64
		mu    sync.Mutex
	)
	data := &ListEntryData{
		header:  listEntryDataHeader,
		entries: make([]*ListEntry, 0, entriesSize),
	}
	err := man.handle(func(_ *Manager, entry *entry) error {
		e := &ListEntry{
			entry: entry,
		}
		mu.Lock()
		data.entries = append(data.entries, e)
		total += e.StoredBytes
		mu.Unlock()
		return nil
	})
	if err != nil {
		return nil, err
	}
	data.TotalStoredBytes = total
	return data, nil
}
