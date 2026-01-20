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
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       -1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
		},
		{
			name: "nil client",
			args: args{
				client: nil,
			},
			want: &Manager{
				client:             nil,
				regions:            DefaultRegions,
				desiredState:       -1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
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
		deletionProtection *bool
		filterExpr         *filterExpr
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
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-west-1", "eu-central-1"},
			},
			wantErr: false,
		},
		{
			name: "empty regions",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{},
			},
			wantErr: false,
		},
		{
			name: "nil regions",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: nil,
			},
			wantErr: false,
		},
		{
			name: "with unsupported region",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-west-1", "invalid-region"},
			},
			wantErr: true,
		},
		{
			name: "with duplicate regions",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"us-west-1", "us-west-1"},
			},
			wantErr: false,
		},
		{
			name: "with uppercase regions",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: []string{"US-WEST-1", "eu-central-1"},
			},
			wantErr: true,
		},
		{
			name: "default regions",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
			},
			args: args{
				regions: DefaultRegions,
			},
			wantErr: false,
		},
		{
			name: "one region",
			fields: fields{
				client:             &Client{},
				regions:            DefaultRegions,
				desiredState:       1,
				deletionProtection: aws.Bool(false),
				filterExpr:         nil,
				sem:                semaphore.NewWeighted(NumWorker),
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
				filterExpr:         tt.fields.filterExpr,
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
		filterExpr         *filterExpr
		sem                *semaphore.Weighted
	}
	type args struct {
		desired string
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
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: DesiredStateOneDay.String(),
			},
			wantErr: false,
		},
		{
			name: "deletion protection",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: DesiredStateProtected.String(),
			},
			wantErr: false,
		},
		{
			name: "invalid desired state",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 0,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				desired: DesiredStateNone.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				client:             tt.fields.client,
				regions:            tt.fields.regions,
				desiredState:       tt.fields.desiredState,
				desiredStateNative: tt.fields.desiredStateNative,
				filterExpr:         tt.fields.filterExpr,
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
		filterExpr         *filterExpr
		filterRaw          string
		sem                *semaphore.Weighted
	}
	type args struct {
		filter string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filter: `name == "error-log"`,
			},
			wantErr: false,
		},
		{
			name: "empty",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filter: ``,
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				client:       &Client{},
				regions:      DefaultRegions,
				desiredState: 1,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(NumWorker),
			},
			args: args{
				filter: `[`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				client:             tt.fields.client,
				regions:            tt.fields.regions,
				desiredState:       tt.fields.desiredState,
				desiredStateNative: tt.fields.desiredStateNative,
				filterExpr:         tt.fields.filterExpr,
				filterRaw:          tt.fields.filterRaw,
				sem:                tt.fields.sem,
			}
			err := man.SetFilter(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_String(t *testing.T) {
	type fields struct {
		regions      []string
		desiredState DesiredState
		filterRaw    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic manager",
			fields: fields{
				regions:      DefaultRegions,
				desiredState: 7,
				filterRaw:    `name == "logname"`,
			},
			want: `{"regions":["ap-northeast-1","ap-northeast-2","ap-northeast-3","ap-south-1","ap-southeast-1","ap-southeast-2","ca-central-1","eu-central-1","eu-west-1","eu-west-2","eu-west-3","eu-north-1","sa-east-1","us-east-1","us-east-2","us-west-1","us-west-2"],"desiredState":"1week","filter":"name == \"logname\""}`,
		},
		{
			name: "empty manager",
			fields: fields{
				regions:      nil,
				desiredState: 0,
				filterRaw:    "",
			},
			want: `{"regions":null,"desiredState":"delete","filter":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				regions:      tt.fields.regions,
				desiredState: tt.fields.desiredState,
				filterRaw:    tt.fields.filterRaw,
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
