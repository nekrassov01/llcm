package llcm

import (
	"sync"
)

// Preview returns the log group entries with the desired state and its simulated results.
func (man *Manager) Preview() (*PreviewEntryData, error) {
	var (
		totalStoredBytes    int64
		totalReducibleBytes int64
		totalRemainingBytes int64
		mu                  sync.Mutex
	)
	data := &PreviewEntryData{
		header:  previewEntryDataHeader,
		entries: make([]*PreviewEntry, 0, entriesSize),
	}
	err := man.handle(func(man *Manager, entry *entry) error {
		e := &PreviewEntry{
			entry: entry,
		}
		e.simulate(man.desiredState)
		mu.Lock()
		data.entries = append(data.entries, e)
		totalStoredBytes += e.StoredBytes
		totalReducibleBytes += e.ReducibleBytes
		totalRemainingBytes += e.RemainingBytes
		mu.Unlock()
		return nil
	})
	if err != nil {
		return nil, err
	}
	data.TotalStoredBytes = totalStoredBytes
	data.TotalReducibleBytes = totalReducibleBytes
	data.TotalRemainingBytes = totalRemainingBytes
	return data, nil
}
