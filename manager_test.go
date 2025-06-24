package llcm

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/semaphore"
)

func TestNewManager(t *testing.T) {
	type args struct {
		client *Client
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "empty client",
			args: args{
				client: &Client{},
			},
			want: &Manager{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: -9999,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
		},
		{
			name: "nil client",
			args: args{
				client: nil,
			},
			want: &Manager{
				client:       nil,
				regions:      DefaultRegions,
				desiredState: -9999,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetRegion(t *testing.T) {
	type fields struct {
		client             *Client
		regions            []string
		desiredState       DesiredState
		desiredStateNative *int32
		filters            []Filter
		filterFns          []func(*entry) bool
		sem                *semaphore.Weighted
	}
	type args struct {
		regions []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid regions",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-west-1", "eu-central-1"},
			},
			wantErr: false,
		},
		{
			name: "empty regions",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{},
			},
			wantErr: false,
		},
		{
			name: "nil regions",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: nil,
			},
			wantErr: false,
		},
		{
			name: "with unsupported region",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-west-1", "invalid-region"},
			},
			wantErr: true,
		},
		{
			name: "with duplicate regions",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-west-1", "us-west-1"},
			},
			wantErr: false,
		},
		{
			name: "with uppercase regions",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"US-WEST-1", "eu-central-1"},
			},
			wantErr: true,
		},
		{
			name: "default regions",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: DefaultRegions,
			},
			wantErr: false,
		},
		{
			name: "one region",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-east-1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				client:             tt.fields.client,
				regions:            tt.fields.regions,
				desiredState:       tt.fields.desiredState,
				desiredStateNative: tt.fields.desiredStateNative,
				filters:            tt.fields.filters,
				filterFns:          tt.fields.filterFns,
				sem:                tt.fields.sem,
			}
			if err := man.SetRegion(tt.args.regions); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetRegion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetDesiredState(t *testing.T) {
	type fields struct {
		client             *Client
		regions            []string
		desiredState       DesiredState
		desiredStateNative *int32
		filters            []Filter
		filterFns          []func(*entry) bool
		sem                *semaphore.Weighted
	}
	type args struct {
		desired DesiredState
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid desired state",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid desired state",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: -9999,
			},
			wantErr: true,
		},
		{
			name: "multiple valid desired states 1",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				desiredStateNative: aws.Int32(1),
				filters:            nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: 2,
			},
			wantErr: false,
		},
		{
			name: "multiple valid desired states 2",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       2,
				desiredStateNative: aws.Int32(2),
				filters:            nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: 2,
			},
			wantErr: false,
		},
		{
			name: "multiple valid desired states 3",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: -9999,
			},
			wantErr: true,
		},
		{
			name: "desired state with nil manager",
			fields: fields{
				client:       nil,
				regions:      nil,
				desiredState: 0,
				filters:      nil,
				sem:          nil,
			},
			args: args{
				desired: 1,
			},
			wantErr: false,
		},
		{
			name: "desired state with max value",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: 2147483647,
			},
			wantErr: false,
		},
		{
			name: "desired state with min value",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: -2147483648,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				client:             tt.fields.client,
				regions:            tt.fields.regions,
				desiredState:       tt.fields.desiredState,
				desiredStateNative: tt.fields.desiredStateNative,
				filters:            tt.fields.filters,
				filterFns:          tt.fields.filterFns,
				sem:                tt.fields.sem,
			}
			if err := man.SetDesiredState(tt.args.desired); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetRetentionInDays() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetFilter(t *testing.T) {
	type fields struct {
		client             *Client
		regions            []string
		desiredState       DesiredState
		desiredStateNative *int32
		filters            []Filter
		filterFns          []func(*entry) bool
		sem                *semaphore.Weighted
	}
	type args struct {
		filters []Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid string filter",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "error-log",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid number filter",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 30,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorGT,
						Value:    "60",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid regex filter",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorREQ,
						Value:    "^error-.*",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid regex filter",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorREQ,
						Value:    "[",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid delete value",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorEQ,
						Value:    "0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid retention value",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 3,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorEQ,
						Value:    "3days", // expected to be converted to 3
					},
				},
			},
			wantErr: false,
		},
		{
			name: "numeric value for string key",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 5,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "12345",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "boolean value for number key",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 10,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorEQ,
						Value:    "true",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty value",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyClass,
						Operator: FilterOperatorEQ,
						Value:    "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "unsupported key type",
			fields: fields{
				client:       &Client{},
				regions:      []string{"us-west-1"},
				desiredState: 10,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKey(12345),
						Operator: FilterOperatorGT,
						Value:    "100",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unsupported operator type",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperator(12345),
						Value:    "value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "zero retention value",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorEQ,
						Value:    "0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil filter",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filters:      nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filters: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				client:             tt.fields.client,
				regions:            tt.fields.regions,
				desiredState:       tt.fields.desiredState,
				desiredStateNative: tt.fields.desiredStateNative,
				filters:            tt.fields.filters,
				filterFns:          tt.fields.filterFns,
				sem:                tt.fields.sem,
			}
			err := man.SetFilter(tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_String(t *testing.T) {
	type fields struct {
		client       *Client
		DesiredState DesiredState
		Filters      []Filter
		Regions      []string
		desiredState *int32
		filterFns    []func(*entry) bool
		sem          *semaphore.Weighted
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic manager",
			fields: fields{
				client:       &Client{},
				DesiredState: 7,
				Filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "logname",
					},
				},
				Regions:      DefaultRegions,
				desiredState: aws.Int32(1),
				filterFns:    nil,
				sem:          nil,
			},
			want: `{
  "regions": [
    "ap-northeast-1",
    "ap-northeast-2",
    "ap-northeast-3",
    "ap-south-1",
    "ap-southeast-1",
    "ap-southeast-2",
    "ca-central-1",
    "eu-central-1",
    "eu-west-1",
    "eu-west-2",
    "eu-west-3",
    "eu-north-1",
    "sa-east-1",
    "us-east-1",
    "us-east-2",
    "us-west-1",
    "us-west-2"
  ],
  "desiredState": "1week",
  "filters": [
    {
      "Key": "name",
      "Operator": "==",
      "Value": "logname"
    }
  ]
}`,
		},
		{
			name: "empty manager",
			fields: fields{
				client:       nil,
				DesiredState: 0,
				Filters:      nil,
				Regions:      nil,
				desiredState: nil,
				filterFns:    nil,
				sem:          nil,
			},
			want: `{
  "regions": null,
  "desiredState": "delete",
  "filters": null
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				client:             tt.fields.client,
				regions:            tt.fields.Regions,
				desiredState:       tt.fields.DesiredState,
				desiredStateNative: tt.fields.desiredState,
				filters:            tt.fields.Filters,
				filterFns:          tt.fields.filterFns,
				sem:                tt.fields.sem,
			}
			if got := man.String(); got != tt.want {
				t.Errorf("Manager.String() = %v, want %v", got, tt.want)
			}
			if diff := cmp.Diff(man.String(), tt.want); diff != "" {
				t.Errorf("Manager.String() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
