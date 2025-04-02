package llcm

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewRenderer(t *testing.T) {
	type args struct {
		data       ListEntryData
		outputType OutputType
	}
	tests := []struct {
		name string
		args args
		want *Renderer[*ListEntry, *ListEntryData]
	}{
		{
			name: "basic",
			args: args{
				data:       listEntryData,
				outputType: OutputTypeJSON,
			},
			want: &Renderer[*ListEntry, *ListEntryData]{
				Data:       &listEntryData,
				OutputType: OutputTypeJSON,
				w:          &bytes.Buffer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got := NewRenderer(w, &tt.args.data, tt.args.outputType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRenderer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderer_String(t *testing.T) {
	type fields struct {
		Data       ListEntryData
		OutputType OutputType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeJSON,
			},
			want: `{
  "Data": {
    "TotalStoredBytes": 0
  },
  "OutputType": "json"
}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := new(bytes.Buffer)
			ren := &Renderer[*ListEntry, *ListEntryData]{
				Data:       &tt.fields.Data,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if got := ren.String(); got != tt.want {
				t.Errorf("Renderer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderer_Render1(t *testing.T) {
	type fields struct {
		Data       ListEntryData
		OutputType OutputType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeJSON,
			},
			want: `[{"LogGroupName":"group0","Region":"ap-northeast-1","Source":"source0","Class":"STANDARD","CreatedAt":"2025-01-01T00:00:00Z","ElapsedDays":90,"RetentionInDays":30,"StoredBytes":1024},{"LogGroupName":"group1","Region":"ap-northeast-2","Source":"source1","Class":"INFREQUENT_ACCESS","CreatedAt":"2024-04-01T00:00:00Z","ElapsedDays":365,"RetentionInDays":30,"StoredBytes":2048}]
`,
			wantErr: false,
		},
		{
			name: "prettyjson",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypePrettyJSON,
			},
			want: `[
  {
    "LogGroupName": "group0",
    "Region": "ap-northeast-1",
    "Source": "source0",
    "Class": "STANDARD",
    "CreatedAt": "2025-01-01T00:00:00Z",
    "ElapsedDays": 90,
    "RetentionInDays": 30,
    "StoredBytes": 1024
  },
  {
    "LogGroupName": "group1",
    "Region": "ap-northeast-2",
    "Source": "source1",
    "Class": "INFREQUENT_ACCESS",
    "CreatedAt": "2024-04-01T00:00:00Z",
    "ElapsedDays": 365,
    "RetentionInDays": 30,
    "StoredBytes": 2048
  }
]
`,
			wantErr: false,
		},
		{
			name: "text",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeText,
			},
			want: `+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
`,
			wantErr: false,
		},
		{
			name: "compressed",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeCompressedText,
			},
			want: `+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+
`,
			wantErr: false,
		},
		{
			name: "markdown",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeMarkdown,
			},
			want: `| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes |
|--------|----------------|---------|-------------------|----------------------|-------------|-----------------|-------------|
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |
`,
			wantErr: false,
		},
		{
			name: "backlog",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeBacklog,
			},
			want: `| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes |h
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |
`,
			wantErr: false,
		},
		{
			name: "tsv",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeTSV,
			},
			want: `Name	Region	Source	Class	CreatedAt	ElapsedDays	RetentionInDays	StoredBytes
group0	ap-northeast-1	source0	STANDARD	2025-01-01T00:00:00Z	90	30	1024
group1	ap-northeast-2	source1	INFREQUENT_ACCESS	2024-04-01T00:00:00Z	365	30	2048
`,
			wantErr: false,
		},
		{
			name: "unknown output type",
			fields: fields{
				Data:       listEntryData,
				OutputType: OutputTypeNone,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "table error",
			fields: fields{
				Data:       errData,
				OutputType: OutputTypeText,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ren := &Renderer[*ListEntry, *ListEntryData]{
				Data:       &tt.fields.Data,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if err := ren.Render(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			opt := cmp.AllowUnexported(ListEntryData{}, ListEntry{}, entry{})
			if diff := cmp.Diff(tt.want, w.String(), opt); diff != "" {
				t.Errorf("Renderer.Render() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRenderer_Render2(t *testing.T) {
	type fields struct {
		Data       PreviewEntryData
		OutputType OutputType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypeJSON,
			},
			want: `[{"LogGroupName":"group0","Region":"ap-northeast-1","Source":"source0","Class":"STANDARD","CreatedAt":"2025-01-01T00:00:00Z","ElapsedDays":90,"RetentionInDays":30,"StoredBytes":1024,"BytesPerDay":0,"DesiredState":0,"ReductionInDays":0,"ReducibleBytes":0,"RemainingBytes":0},{"LogGroupName":"group1","Region":"ap-northeast-2","Source":"source1","Class":"INFREQUENT_ACCESS","CreatedAt":"2024-04-01T00:00:00Z","ElapsedDays":365,"RetentionInDays":30,"StoredBytes":2048,"BytesPerDay":100,"DesiredState":100,"ReductionInDays":100,"ReducibleBytes":100,"RemainingBytes":100}]
`,
			wantErr: false,
		},
		{
			name: "prettyjson",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypePrettyJSON,
			},
			want: `[
  {
    "LogGroupName": "group0",
    "Region": "ap-northeast-1",
    "Source": "source0",
    "Class": "STANDARD",
    "CreatedAt": "2025-01-01T00:00:00Z",
    "ElapsedDays": 90,
    "RetentionInDays": 30,
    "StoredBytes": 1024,
    "BytesPerDay": 0,
    "DesiredState": 0,
    "ReductionInDays": 0,
    "ReducibleBytes": 0,
    "RemainingBytes": 0
  },
  {
    "LogGroupName": "group1",
    "Region": "ap-northeast-2",
    "Source": "source1",
    "Class": "INFREQUENT_ACCESS",
    "CreatedAt": "2024-04-01T00:00:00Z",
    "ElapsedDays": 365,
    "RetentionInDays": 30,
    "StoredBytes": 2048,
    "BytesPerDay": 100,
    "DesiredState": 100,
    "ReductionInDays": 100,
    "ReducibleBytes": 100,
    "RemainingBytes": 100
  }
]
`,
			wantErr: false,
		},
		{
			name: "text",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypeText,
			},
			want: `+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes | BytesPerDay | DesiredState | ReductionInDays | ReducibleBytes | RemainingBytes |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |           0 |            0 |               0 |              0 |              0 |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |         100 |          100 |             100 |            100 |            100 |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
`,
			wantErr: false,
		},
		{
			name: "compressed",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypeCompressedText,
			},
			want: `+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes | BytesPerDay | DesiredState | ReductionInDays | ReducibleBytes | RemainingBytes |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |           0 |            0 |               0 |              0 |              0 |
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |         100 |          100 |             100 |            100 |            100 |
+--------+----------------+---------+-------------------+----------------------+-------------+-----------------+-------------+-------------+--------------+-----------------+----------------+----------------+
`,
			wantErr: false,
		},
		{
			name: "markdown",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypeMarkdown,
			},
			want: `| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes | BytesPerDay | DesiredState | ReductionInDays | ReducibleBytes | RemainingBytes |
|--------|----------------|---------|-------------------|----------------------|-------------|-----------------|-------------|-------------|--------------|-----------------|----------------|----------------|
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |           0 |            0 |               0 |              0 |              0 |
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |         100 |          100 |             100 |            100 |            100 |
`,
			wantErr: false,
		},
		{
			name: "backlog",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypeBacklog,
			},
			want: `| Name   | Region         | Source  | Class             | CreatedAt            | ElapsedDays | RetentionInDays | StoredBytes | BytesPerDay | DesiredState | ReductionInDays | ReducibleBytes | RemainingBytes |h
| group0 | ap-northeast-1 | source0 | STANDARD          | 2025-01-01T00:00:00Z |          90 |              30 |        1024 |           0 |            0 |               0 |              0 |              0 |
| group1 | ap-northeast-2 | source1 | INFREQUENT_ACCESS | 2024-04-01T00:00:00Z |         365 |              30 |        2048 |         100 |          100 |             100 |            100 |            100 |
`,
			wantErr: false,
		},
		{
			name: "tsv",
			fields: fields{
				Data:       previewEntryData,
				OutputType: OutputTypeTSV,
			},
			want: `Name	Region	Source	Class	CreatedAt	ElapsedDays	RetentionInDays	StoredBytes	BytesPerDay	DesiredState	ReductionInDays	ReducibleBytes	RemainingBytes
group0	ap-northeast-1	source0	STANDARD	2025-01-01T00:00:00Z	90	30	1024	0	0	0	0	0
group1	ap-northeast-2	source1	INFREQUENT_ACCESS	2024-04-01T00:00:00Z	365	30	2048	100	100	100	100	100
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ren := &Renderer[*PreviewEntry, *PreviewEntryData]{
				Data:       &tt.fields.Data,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if err := ren.Render(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := w.String(); got != tt.want {
				t.Errorf("Renderer.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderer_toChart1(t *testing.T) {
	type fields struct {
		Data       *ListEntryData
		OutputType OutputType
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "list",
			fields: fields{
				Data:       &listEntryData,
				OutputType: OutputTypeChart,
			},
			wantErr: false,
		},
		{
			name: "nil",
			fields: fields{
				Data:       &ListEntryData{},
				OutputType: OutputTypeChart,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ren := &Renderer[*ListEntry, *ListEntryData]{
				Data:       tt.fields.Data,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if err := ren.toChart(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.toChart() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := ren.toChart(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.toChart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenderer_toChart2(t *testing.T) {
	type fields struct {
		Data       *PreviewEntryData
		OutputType OutputType
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "list",
			fields: fields{
				Data:       &previewEntryData,
				OutputType: OutputTypeChart,
			},
			wantErr: false,
		},
		{
			name: "nil",
			fields: fields{
				Data:       &PreviewEntryData{},
				OutputType: OutputTypeChart,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ren := &Renderer[*PreviewEntry, *PreviewEntryData]{
				Data:       tt.fields.Data,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if err := ren.toChart(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.toChart() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := ren.toChart(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.toChart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
