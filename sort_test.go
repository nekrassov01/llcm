package llcm

import (
	"reflect"
	"testing"
)

func TestSortEntries(t *testing.T) {
	data := &ListEntryData{
		entries: []*ListEntry{
			{
				entry: &entry{
					LogGroupName: "2",
					StoredBytes:  200,
				},
			},
			{
				entry: &entry{
					LogGroupName: "5",
					StoredBytes:  100,
				},
			},
			{
				entry: &entry{
					LogGroupName: "4",
					StoredBytes:  100,
				},
			},
			{
				entry: &entry{
					LogGroupName: "1",
					StoredBytes:  300,
				},
			},
			{
				entry: &entry{
					LogGroupName: "3",
					StoredBytes:  150,
				},
			},
		},
	}
	sorted := &ListEntryData{
		entries: []*ListEntry{
			{
				entry: &entry{
					LogGroupName: "1",
					StoredBytes:  300,
				},
			},
			{
				entry: &entry{
					LogGroupName: "2",
					StoredBytes:  200,
				},
			},
			{
				entry: &entry{
					LogGroupName: "3",
					StoredBytes:  150,
				},
			},
			{
				entry: &entry{
					LogGroupName: "4",
					StoredBytes:  100,
				},
			},
			{
				entry: &entry{
					LogGroupName: "5",
					StoredBytes:  100,
				},
			},
		},
	}
	SortEntries(data)
	for i, entry := range data.entries {
		var (
			got  = entry
			want = sorted.entries[i]
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Entry[%d] = %v, want %v", i, got, want)
		}
	}
}
