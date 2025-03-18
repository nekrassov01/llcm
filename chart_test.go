package llcm

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Test_getPieItems(t *testing.T) {
	type args struct {
		entries []*ListEntry
	}
	type want struct {
		title string
		items []opts.PieData
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "basic",
			args: args{
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     1024,
						},
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     4096,
						},
					},
				},
			},
			want: want{
				title: "Stored bytes of log groups",
				items: []opts.PieData{
					{
						Name:  "group0",
						Value: int64(1024),
					},
					{
						Name:  "group1",
						Value: int64(4096),
					},
				},
			},
		},
		{
			name: "include zero",
			args: args{
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     0,
						},
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     4096,
						},
					},
				},
			},
			want: want{
				title: "Stored bytes of log groups",
				items: []opts.PieData{
					{
						Name:  "group1",
						Value: int64(4096),
					},
				},
			},
		},
		{
			name: "others",
			args: args{
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     1024,
						},
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     4096,
						},
					},
					{
						entry: &entry{
							LogGroupName:    "group2",
							Region:          "us-east-1",
							Source:          "source2",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     256,
						},
					},
					{
						entry: &entry{
							LogGroupName:    "group3",
							Region:          "us-east-1",
							Source:          "source3",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     512,
						},
					},
				},
			},
			want: want{
				title: "Stored bytes of log groups",
				items: []opts.PieData{
					{
						Name:  "group0",
						Value: int64(1024),
					},
					{
						Name:  "group1",
						Value: int64(4096),
					},
					{
						Name:  "others",
						Value: int64(768),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, items := getPieItems(tt.args.entries)
			if title != tt.want.title {
				t.Errorf("getPieItems() title = %v, want %v", title, tt.want)
			}
			if !reflect.DeepEqual(items, tt.want.items) {
				t.Errorf("getPieItems() items = %v, want %v", items, tt.want.items)
			}
		})
	}
}

func Test_renderPieChart(t *testing.T) {
	type args struct {
		pie *charts.Pie
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				pie: newPieChart("Stored bytes of log groups", []opts.PieData{
					{
						Name:  "group1",
						Value: float64(8192),
					},
					{
						Name:  "group0",
						Value: float64(2048),
					},
					{
						Name:  "others",
						Value: float64(512),
					},
				}),
			},
			wantErr: false,
		},
		{
			name: "nil",
			args: args{
				pie: newPieChart("Stored bytes of log groups", nil),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := renderPieChart(tt.args.pie); (err != nil) != tt.wantErr {
				t.Errorf("renderPieChart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getBarTitle(t *testing.T) {
	type args struct {
		entries []*PreviewEntry
	}
	type want struct {
		title    string
		subtitle string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "desired state 0",
			args: args{
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 9999,
							StoredBytes:     1024,
						},
						BytesPerDay:     1024,
						DesiredState:    0,
						ReductionInDays: 1,
						ReducibleBytes:  0,
						RemainingBytes:  1024,
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 731,
							StoredBytes:     0,
						},
						BytesPerDay:     0,
						DesiredState:    0,
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  0,
					},
				},
			},
			want: want{
				title:    "The simulation of reductions in log groups",
				subtitle: "Desired state: Delete log groups",
			},
		},
		{
			name: "desired state 9999",
			args: args{
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     1024,
						},
						BytesPerDay:     1024,
						DesiredState:    9999,
						ReductionInDays: 1,
						ReducibleBytes:  0,
						RemainingBytes:  1024,
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 731,
							StoredBytes:     0,
						},
						BytesPerDay:     0,
						DesiredState:    9999,
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  0,
					},
				},
			},
			want: want{
				title:    "The simulation of reductions in log groups",
				subtitle: "Desired state: Delete retention policy",
			},
		},
		{
			name: "desired state 365",
			args: args{
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 30,
							StoredBytes:     1024,
						},
						BytesPerDay:     1024,
						DesiredState:    365,
						ReductionInDays: 1,
						ReducibleBytes:  0,
						RemainingBytes:  1024,
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     0,
							RetentionInDays: 731,
							StoredBytes:     0,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  0,
					},
				},
			},
			want: want{
				title:    "The simulation of reductions in log groups",
				subtitle: "Desired state: Change retention to 365 days",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, subtitle := getBarTitle(tt.args.entries)
			if title != tt.want.title {
				t.Errorf("getBarTitle() got = %v, title %v", title, tt.want.title)
			}
			if subtitle != tt.want.subtitle {
				t.Errorf("getBarTitle() subtitle = %v, subtitle %v", subtitle, tt.want.subtitle)
			}
		})
	}
}

func Test_getBarItems(t *testing.T) {
	type args struct {
		entries []*PreviewEntry
	}
	type want struct {
		names      []string
		remainings []opts.BarData
		reducibles []opts.BarData
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "basic",
			args: args{
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     6144,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  5120,
						RemainingBytes:  1024,
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     3072,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  512,
						RemainingBytes:  2560,
					},
				},
			},
			want: want{
				names: []string{"group0", "group1"},
				remainings: []opts.BarData{
					{
						Value: int64(1024),
					},
					{
						Value: int64(2560),
					},
				},
				reducibles: []opts.BarData{
					{
						Value: int64(5120),
					},
					{
						Value: int64(512),
					},
				},
			},
		},
		{
			name: "include zero",
			args: args{
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     0,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  0,
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     3072,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  512,
						RemainingBytes:  2560,
					},
				},
			},
			want: want{
				names: []string{"group1"},
				remainings: []opts.BarData{
					{
						Value: int64(2560),
					},
				},
				reducibles: []opts.BarData{
					{
						Value: int64(512),
					},
				},
			},
		},
		{
			name: "others",
			args: args{
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "group0",
							Region:          "ap-northeast-1",
							Source:          "source0",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     6144,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  5120,
						RemainingBytes:  1024,
					},
					{
						entry: &entry{
							LogGroupName:    "group1",
							Region:          "ap-northeast-2",
							Source:          "source1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     3072,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  512,
						RemainingBytes:  2560,
					},
					{
						entry: &entry{
							LogGroupName:    "group2",
							Region:          "us-east-1",
							Source:          "source2",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     256,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  512,
						RemainingBytes:  256,
					},
					{
						entry: &entry{
							LogGroupName:    "group3",
							Region:          "us-east-1",
							Source:          "source3",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       time.Now(),
							ElapsedDays:     3000,
							RetentionInDays: 2192,
							StoredBytes:     512,
						},
						BytesPerDay:     0,
						DesiredState:    365,
						ReductionInDays: 0,
						ReducibleBytes:  512,
						RemainingBytes:  0,
					},
				},
			},
			want: want{
				names: []string{"group0", "group1", "others"},
				remainings: []opts.BarData{
					{
						Value: int64(1024),
					},
					{
						Value: int64(2560),
					},
					{
						Value: int64(256),
					},
				},
				reducibles: []opts.BarData{
					{
						Value: int64(5120),
					},
					{
						Value: int64(512),
					},
					{
						Value: int64(1024),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			names, remainings, reducibles := getBarItems(tt.args.entries)
			if !reflect.DeepEqual(names, tt.want.names) {
				t.Errorf("getBarItems() names = %v, want %v", names, tt.want.names)
			}
			if !reflect.DeepEqual(remainings, tt.want.remainings) {
				t.Errorf("getBarItems() remainings = %v, want %v", remainings, tt.want.remainings)
			}
			if !reflect.DeepEqual(reducibles, tt.want.reducibles) {
				t.Errorf("getBarItems() reducibles = %v, want %v", reducibles, tt.want.reducibles)
			}
		})
	}
}

func Test_renderBarChart(t *testing.T) {
	type args struct {
		bar *charts.Bar
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete log groups",
			args: args{
				bar: newBarChart(
					"The simulation of reductions in log groups",
					"Desired state: Change retention to 365 days",
					[]string{"/aws/lambda/loggroup1-0123456789abcdef", "/aws/lambda/loggroup2-0123456789abcdef", "/aws/lambda/loggroup0-0123456789abcdef", "others"},
					[]opts.BarData{
						{
							Name:  "/aws/lambda/loggroup1-0123456789abcdef",
							Value: float64(2048),
						},
						{
							Name:  "/aws/lambda/loggroup2-0123456789abcdef",
							Value: float64(1024),
						},
						{
							Name:  "/aws/lambda/loggroup0-0123456789abcdef",
							Value: float64(2048),
						},
						{
							Name:  "others",
							Value: float64(128),
						},
					},
					[]opts.BarData{
						{
							Name:  "/aws/lambda/loggroup1-0123456789abcdef",
							Value: float64(6144),
						},
						{
							Name:  "/aws/lambda/loggroup2-0123456789abcdef",
							Value: float64(3072),
						},
						{
							Name:  "/aws/lambda/loggroup0-0123456789abcdef",
							Value: float64(0),
						},
						{
							Name:  "others",
							Value: float64(384),
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "nil",
			args: args{
				bar: newBarChart("The simulation of reductions in log groups", "", nil, nil, nil),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := renderBarChart(tt.args.bar); (err != nil) != tt.wantErr {
				t.Errorf("renderBarChart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
