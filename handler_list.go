package llcm

import (
	"context"
	"sync"
)

// List lists the log group entries.
func (man *Manager) List(ctx context.Context) (*ListEntryData, error) {
	var (
		total int64
		mu    sync.Mutex
	)
	data := &ListEntryData{
		header:  listEntryDataHeader,
		entries: make([]*ListEntry, 0, entriesSize),
	}
	fn := func(entry *entry) error {
		e := &ListEntry{
			entry: entry,
		}
		mu.Lock()
		data.entries = append(data.entries, e)
		total += e.StoredBytes
		mu.Unlock()
		return nil
	}
	if err := man.handle(ctx, fn); err != nil {
		return nil, err
	}
	data.TotalStoredBytes = total
	return data, nil
}
