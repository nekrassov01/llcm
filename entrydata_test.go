package llcm

import (
	"reflect"
	"testing"
)

func TestListEntryData_Total(t *testing.T) {
	type fields struct {
		TotalStoredBytes int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int64
	}{
		{
			name: "basic",
			fields: fields{
				TotalStoredBytes: 100,
			},
			want: map[string]int64{
				TotalStoredBytesLabel: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ListEntryData{
				TotalStoredBytes: tt.fields.TotalStoredBytes,
			}
			if got := d.Total(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListEntryData.Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListEntryData_Chart(t *testing.T) {
	type fields struct {
		entries []*ListEntry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				entries: listEntryData.entries,
			},
			wantErr: false,
		},
		{
			name: "nil",
			fields: fields{
				entries: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ListEntryData{
				entries: tt.fields.entries,
			}
			if err := d.Chart(); (err != nil) != tt.wantErr {
				t.Errorf("ListEntryData.Chart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreviewEntryData_Total(t *testing.T) {
	type fields struct {
		TotalStoredBytes    int64
		TotalReducibleBytes int64
		TotalRemainingBytes int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int64
	}{
		{
			name: "basic",
			fields: fields{
				TotalStoredBytes:    100,
				TotalReducibleBytes: 50,
				TotalRemainingBytes: 50,
			},
			want: map[string]int64{
				TotalStoredBytesLabel:    100,
				TotalReducibleBytesLabel: 50,
				TotalRemainingBytesLabel: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &PreviewEntryData{
				TotalStoredBytes:    tt.fields.TotalStoredBytes,
				TotalReducibleBytes: tt.fields.TotalReducibleBytes,
				TotalRemainingBytes: tt.fields.TotalRemainingBytes,
			}
			if got := d.Total(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreviewEntryData.Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreviewEntryData_Chart(t *testing.T) {
	type fields struct {
		TotalStoredBytes    int64
		TotalReducibleBytes int64
		TotalRemainingBytes int64
		header              []string
		entries             []*PreviewEntry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				entries: previewEntryData.entries,
			},
			wantErr: false,
		},
		{
			name: "nil",
			fields: fields{
				entries: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &PreviewEntryData{
				TotalStoredBytes:    tt.fields.TotalStoredBytes,
				TotalReducibleBytes: tt.fields.TotalReducibleBytes,
				TotalRemainingBytes: tt.fields.TotalRemainingBytes,
				header:              tt.fields.header,
				entries:             tt.fields.entries,
			}
			if err := d.Chart(); (err != nil) != tt.wantErr {
				t.Errorf("PreviewEntryData.Chart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
