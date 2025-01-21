package llcm

import (
	"cmp"
	"slices"
)

// SortEntries sorts the entries by bytes and name.
func SortEntries[T Entry](entries []T) {
	slices.SortFunc(entries, func(a, b T) int {
		if n := cmp.Compare(b.Bytes(), a.Bytes()); n != 0 {
			return n
		}
		return cmp.Compare(a.Name(), b.Name())
	})
}
