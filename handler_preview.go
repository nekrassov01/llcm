package llcm

import (
	"sync"
	"sync/atomic"
)

// Preview returns the log group entries with the desired state and its simulated results.
func (man *Manager) Preview() (*PreviewEntryData, error) {
	var (
		totalStoredBytes    int64
		totalReducibleBytes int64
		totalRemainingBytes int64
		wg                  sync.WaitGroup
		entryChan           = make(chan *PreviewEntry, regionalEntriesSize)
		data                = &PreviewEntryData{
			header:  previewEntryDataHeader,
			entries: make([]*PreviewEntry, 0, globalEntriesSize),
		}
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for entry := range entryChan {
			data.entries = append(data.entries, entry)
			atomic.AddInt64(&totalStoredBytes, entry.StoredBytes)
			atomic.AddInt64(&totalReducibleBytes, entry.ReducibleBytes)
			atomic.AddInt64(&totalRemainingBytes, entry.RemainingBytes)
		}
	}()
	err := man.handle(func(man *Manager, entry *entry) error {
		e := &PreviewEntry{
			entry: entry,
		}
		e.simulate(man.desiredState)
		select {
		case entryChan <- e:
			return nil
		case <-man.ctx.Done():
			return man.ctx.Err()
		}
	})
	close(entryChan)
	if err != nil {
		return nil, err
	}
	wg.Wait()
	data.TotalStoredBytes = atomic.LoadInt64(&totalStoredBytes)
	data.TotalReducibleBytes = atomic.LoadInt64(&totalReducibleBytes)
	data.TotalRemainingBytes = atomic.LoadInt64(&totalRemainingBytes)
	return data, nil
}
