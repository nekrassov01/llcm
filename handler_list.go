package llcm

import (
	"sync"
	"sync/atomic"
)

// List lists the log group entries.
func (man *Manager) List() (*ListEntryData, error) {
	var (
		total     int64
		wg        sync.WaitGroup
		entryChan = make(chan *ListEntry, regionalEntriesSize)
		data      = &ListEntryData{
			header:  listEntryDataHeader,
			entries: make([]*ListEntry, 0, globalEntriesSize),
		}
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for entry := range entryChan {
			data.entries = append(data.entries, entry)
			atomic.AddInt64(&total, entry.StoredBytes)
		}
	}()
	err := man.handle(func(man *Manager, entry *entry) error {
		e := &ListEntry{
			entry: entry,
		}
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
	data.TotalStoredBytes = atomic.LoadInt64(&total)
	return data, nil
}
