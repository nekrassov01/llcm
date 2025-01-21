package llcm

import (
	"reflect"
	"testing"
)

func TestSortEntries(t *testing.T) {
	entries := []Entry{
		&ListEntry{
			entry: &entry{
				LogGroupName: "2",
				StoredBytes:  200,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "5",
				StoredBytes:  100,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "4",
				StoredBytes:  100,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "1",
				StoredBytes:  300,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "3",
				StoredBytes:  150,
			},
		},
	}
	sorted := []Entry{
		&ListEntry{
			entry: &entry{
				LogGroupName: "1",
				StoredBytes:  300,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "2",
				StoredBytes:  200,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "3",
				StoredBytes:  150,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "4",
				StoredBytes:  100,
			},
		},
		&ListEntry{
			entry: &entry{
				LogGroupName: "5",
				StoredBytes:  100,
			},
		},
	}
	SortEntries(entries)
	for i, entry := range entries {
		var (
			got  = entry.(*ListEntry)
			want = sorted[i].(*ListEntry)
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Entry[%d] = %v, want %v", i, got, want)
		}
	}
}
