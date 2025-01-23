package llcm

import (
	"cmp"
	"slices"
)

// SortEntries sorts the entries by bytes and name.
func SortEntries[E Entry, D EntryData[E]](data D) {
	slices.SortFunc(data.Entries(), func(a, b E) int {
		if n := cmp.Compare(b.Bytes(), a.Bytes()); n != 0 {
			return n
		}
		return cmp.Compare(a.Name(), b.Name())
	})
}
