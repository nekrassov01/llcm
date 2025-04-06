package llcm

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"golang.org/x/sync/semaphore"
)

func TestManager_Apply(t *testing.T) {
	type fields struct {
		client             *Client
		regions            []string
		desiredState       DesiredState
		desiredStateNative *int32
		filters            []Filter
		filterFns          []func(*entry) bool
		sem                *semaphore.Weighted
		ctx                context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    int32
		wantErr bool
	}{
		{
			name: "none",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateNone,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "delete log group",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
					DeleteLogGroupFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteLogGroupInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error) {
						return &cloudwatchlogs.DeleteLogGroupOutput{}, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "delete retention policy",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
					DeleteRetentionPolicyFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error) {
						return &cloudwatchlogs.DeleteRetentionPolicyOutput{}, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateInfinite,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "put retention policy",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
					PutRetentionPolicyFunc: func(_ context.Context, _ *cloudwatchlogs.PutRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
						return &cloudwatchlogs.PutRetentionPolicyOutput{}, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneDay,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "delete log group returns error",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
					DeleteLogGroupFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteLogGroupInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error) {
						return nil, errors.New("error")
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "delete retention policy returns error",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
					DeleteRetentionPolicyFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error) {
						return nil, errors.New("error")
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateInfinite,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "put retention policy",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
					PutRetentionPolicyFunc: func(_ context.Context, _ *cloudwatchlogs.PutRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
						return nil, errors.New("error")
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneDay,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "multiple log groups",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group-1"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group-1"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
								{
									LogGroupName:    aws.String("test-log-group-2"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group-2"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
					DeleteLogGroupFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteLogGroupInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error) {
						return &cloudwatchlogs.DeleteLogGroupOutput{}, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filters:      nil,
				sem:          semaphore.NewWeighted(10),
				ctx:          context.Background(),
			},
			want:    2,
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
				ctx:                tt.fields.ctx,
			}
			got, err := man.Apply(io.Discard)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
