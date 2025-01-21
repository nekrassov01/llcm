package llcm

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/semaphore"
)

func TestNewManager(t *testing.T) {
	type args struct {
		ctx    context.Context
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
				ctx:    context.Background(),
				client: &Client{},
			},
			want: &Manager{
				Client:       &Client{},
				DesiredState: -9999,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
		},
		{
			name: "nil client",
			args: args{
				ctx:    context.Background(),
				client: nil,
			},
			want: &Manager{
				Client:       nil,
				DesiredState: -9999,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.ctx, tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetRegions(t *testing.T) {
	type fields struct {
		Client       *Client
		DesiredState DesiredState
		Filters      []Filter
		Regions      []string
		desiredState *int32
		filterFns    []func(*entry) bool
		sem          *semaphore.Weighted
		ctx          context.Context
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
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: []string{"us-west-1", "eu-central-1"},
			},
			wantErr: false,
		},
		{
			name: "empty regions",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: []string{},
			},
			wantErr: false,
		},
		{
			name: "nil regions",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: nil,
			},
			wantErr: false,
		},
		{
			name: "with unsupported region",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: []string{"us-west-1", "invalid-region"},
			},
			wantErr: true,
		},
		{
			name: "with duplicate regions",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: []string{"us-west-1", "us-west-1"},
			},
			wantErr: false,
		},
		{
			name: "with uppercase regions",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: []string{"US-WEST-1", "eu-central-1"},
			},
			wantErr: true,
		},
		{
			name: "default regions",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				regions: DefaultRegions,
			},
			wantErr: false,
		},
		{
			name: "one region",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       tt.fields.Client,
				DesiredState: tt.fields.DesiredState,
				Filters:      tt.fields.Filters,
				Regions:      tt.fields.Regions,
				desiredState: tt.fields.desiredState,
				filterFns:    tt.fields.filterFns,
				sem:          tt.fields.sem,
				ctx:          tt.fields.ctx,
			}
			if err := man.SetRegions(tt.args.regions); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetRegions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetDesiredState(t *testing.T) {
	type fields struct {
		Client       *Client
		DesiredState DesiredState
		Filters      []Filter
		Regions      []string
		desiredState *int32
		filterFns    []func(*entry) bool
		sem          *semaphore.Weighted
		ctx          context.Context
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
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				desired: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid desired state",
			fields: fields{
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				desired: -9999,
			},
			wantErr: true,
		},
		{
			name: "multiple valid desired states 1",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				desiredState: aws.Int32(1),
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				desired: 2,
			},
			wantErr: false,
		},
		{
			name: "multiple valid desired states 2",
			fields: fields{
				Client:       &Client{},
				DesiredState: 2,
				Filters:      nil,
				Regions:      DefaultRegions,
				desiredState: aws.Int32(2),
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				desired: 2,
			},
			wantErr: false,
		},
		{
			name: "multiple valid desired states 3",
			fields: fields{
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				desired: -9999,
			},
			wantErr: true,
		},
		{
			name: "desired state with nil manager",
			fields: fields{
				Client:       nil,
				DesiredState: 0,
				Filters:      nil,
				Regions:      nil,
				sem:          nil,
				ctx:          context.Background(),
			},
			args: args{
				desired: 1,
			},
			wantErr: false,
		},
		{
			name: "desired state with max value",
			fields: fields{
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
			},
			args: args{
				desired: 2147483647,
			},
			wantErr: false,
		},
		{
			name: "desired state with min value",
			fields: fields{
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       tt.fields.Client,
				DesiredState: tt.fields.DesiredState,
				Filters:      tt.fields.Filters,
				Regions:      tt.fields.Regions,
				desiredState: tt.fields.desiredState,
				filterFns:    tt.fields.filterFns,
				sem:          tt.fields.sem,
				ctx:          tt.fields.ctx,
			}
			if err := man.SetDesiredState(tt.args.desired); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetRetentionInDays() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetFilter(t *testing.T) {
	type fields struct {
		Client       *Client
		DesiredState DesiredState
		Filters      []Filter
		Regions      []string
		desiredState *int32
		filterFns    []func(*entry) bool
		sem          *semaphore.Weighted
		ctx          context.Context
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
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 30,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 3,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 5,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 10,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 10,
				Filters:      nil,
				Regions:      []string{"us-west-1"},
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 1,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       &Client{},
				DesiredState: 0,
				Filters:      nil,
				Regions:      DefaultRegions,
				sem:          semaphore.NewWeighted(NumWorker),
				ctx:          context.Background(),
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
				Client:       tt.fields.Client,
				DesiredState: tt.fields.DesiredState,
				Filters:      tt.fields.Filters,
				Regions:      tt.fields.Regions,
				desiredState: tt.fields.desiredState,
				filterFns:    tt.fields.filterFns,
				sem:          tt.fields.sem,
				ctx:          tt.fields.ctx,
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
		Client       *Client
		DesiredState DesiredState
		Filters      []Filter
		Regions      []string
		desiredState *int32
		filterFns    []func(*entry) bool
		sem          *semaphore.Weighted
		ctx          context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal manager",
			fields: fields{
				Client:       &Client{},
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
				ctx:          context.Background(),
			},
			want: `{
  "DesiredState": "1week",
  "Filters": [
    {
      "Key": "name",
      "Operator": "==",
      "Value": "logname"
    }
  ],
  "Regions": [
    "us-east-1",
    "us-east-2",
    "us-west-1",
    "us-west-2",
    "ap-south-1",
    "ap-northeast-3",
    "ap-northeast-2",
    "ap-southeast-1",
    "ap-southeast-2",
    "ap-northeast-1",
    "ca-central-1",
    "eu-central-1",
    "eu-west-1",
    "eu-west-2",
    "eu-west-3",
    "eu-north-1",
    "sa-east-1"
  ]
}`,
		},
		{
			name: "empty manager",
			fields: fields{
				Client:       nil,
				DesiredState: 0,
				Filters:      nil,
				Regions:      nil,
				desiredState: nil,
				filterFns:    nil,
				sem:          nil,
				ctx:          context.Background(),
			},
			want: `{
  "DesiredState": "delete",
  "Filters": null,
  "Regions": null
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:       tt.fields.Client,
				DesiredState: tt.fields.DesiredState,
				Filters:      tt.fields.Filters,
				Regions:      tt.fields.Regions,
				desiredState: tt.fields.desiredState,
				filterFns:    tt.fields.filterFns,
				sem:          tt.fields.sem,
				ctx:          tt.fields.ctx,
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
