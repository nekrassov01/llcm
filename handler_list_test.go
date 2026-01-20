package llcm

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/google/go-cmp/cmp"
	"github.com/nekrassov01/filter"
	"golang.org/x/sync/semaphore"
)

func TestManager_List(t *testing.T) {
	type fields struct {
		client             *Client
		regions            []string
		desiredState       DesiredState
		desiredStateNative *int32
		filterExpr         *filterExpr
		sem                *semaphore.Weighted
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListEntryData
		wantErr bool
	}{
		{
			name: "no entries",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header:  listEntryDataHeader,
				entries: make([]*ListEntry, 0, entriesSize),
			},
			wantErr: false,
		},
		{
			name: "nil entries",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: nil,
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header:  listEntryDataHeader,
				entries: make([]*ListEntry, 0, entriesSize),
			},
			wantErr: false,
		},
		{
			name: "empty result",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						return &cloudwatchlogs.DescribeLogGroupsOutput{}, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header:  listEntryDataHeader,
				entries: make([]*ListEntry, 0, entriesSize),
			},
			wantErr: false,
		},
		{
			name: "single entry",
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
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group"),
						},
					},
				},
				TotalStoredBytes: 1024,
			},
			wantErr: false,
		},
		{
			name: "multiple entries",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
				},
				TotalStoredBytes: 3072,
			},
			wantErr: false,
		},
		{
			name: "pagination",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						if params.NextToken == nil {
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
										LogGroupClass:   types.LogGroupClassInfrequentAccess,
										CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
										RetentionInDays: aws.Int32(7),
										StoredBytes:     aws.Int64(2048),
									},
								},
								NextToken: aws.String("token0"),
							}
							return out, nil
						}
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group-3"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-2:000000000000:log-group:test-log-group-3"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(3),
									StoredBytes:     aws.Int64(0),
								},
								{
									LogGroupName:    aws.String("test-log-group-4"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-2:111111111111:log-group:test-log-group-4"),
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(731),
									StoredBytes:     aws.Int64(4096),
								},
							},
							NextToken: nil,
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-4",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 731,
							StoredBytes:     4096,
							name:            aws.String("test-log-group-4"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-3",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 3,
							StoredBytes:     0,
							name:            aws.String("test-log-group-3"),
						},
					},
				},
				TotalStoredBytes: 7168,
			},
			wantErr: false,
		},
		{
			name: "multiple regions",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
				},
				TotalStoredBytes: 6144,
			},
			wantErr: false,
		},
		{
			name: "zero retention",
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
									RetentionInDays: aws.Int32(0),
									StoredBytes:     aws.Int64(1024),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 9999,
							StoredBytes:     1024,
							name:            aws.String("test-log-group"),
						},
					},
				},
				TotalStoredBytes: 1024,
			},
			wantErr: false,
		},
		{
			name: "with filter name",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`name == "test-log-group-1"`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
				},
				TotalStoredBytes: 2048,
			},
			wantErr: false,
		},
		{
			name: "with filter class",
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
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:000000000000:log-group:test-log-group-2"),
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`class ==* "standard"`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 365,
							StoredBytes:     1024,
							name:            aws.String("test-log-group-1"),
						},
					},
				},
				TotalStoredBytes: 2048,
			},
			wantErr: false,
		},
		{
			name: "deletion protection enabled",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:              aws.String("test-log-group"),
									LogGroupArn:               aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:             types.LogGroupClassStandard,
									CreationTime:              aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays:           aws.Int32(365),
									StoredBytes:               aws.Int64(1024),
									DeletionProtectionEnabled: aws.Bool(true),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateProtected,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`protected == true`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:       "test-log-group",
							Region:             "us-east-1",
							Class:              types.LogGroupClassStandard,
							CreatedAt:          mustTime("2025-01-01T00:00:00Z"),
							DeletionProtection: true,
							ElapsedDays:        90,
							RetentionInDays:    365,
							StoredBytes:        1024,
							name:               aws.String("test-log-group"),
						},
					},
				},
				TotalStoredBytes: 1024,
			},
			wantErr: false,
		},
		{
			name: "with filter elepased days",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group-1"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group-1"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-10T00:00:00Z")),
									RetentionInDays: aws.Int32(365),
									StoredBytes:     aws.Int64(1024),
								},
								{
									LogGroupName:    aws.String("test-log-group-2"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group-2"),
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`elapsed >= 90`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
				},
				TotalStoredBytes: 4096,
			},
			wantErr: false,
		},
		{
			name: "with filter retention 1",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`retention == 7`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
				},
				TotalStoredBytes: 4096,
			},
			wantErr: false,
		},
		{
			name: "with filter retention 2",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr: func() *filterExpr {
					expr, _ := filter.Parse(fmt.Sprintf(`retention == %d`, func() int32 { n, _ := ParseDesiredState("1week"); return int32(n) }()))
					return expr
				}(),
				sem: semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
				},
				TotalStoredBytes: 4096,
			},
			wantErr: false,
		},
		{
			name: "with filter bytes",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`bytes >= 2048`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &ListEntryData{
				header: listEntryDataHeader,
				entries: []*ListEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassInfrequentAccess,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: 7,
							StoredBytes:     2048,
							name:            aws.String("test-log-group-2"),
						},
					},
				},
				TotalStoredBytes: 4096,
			},
			wantErr: false,
		},
		{
			name: "with filter none",
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
									LogGroupClass:   types.LogGroupClassInfrequentAccess,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(7),
									StoredBytes:     aws.Int64(2048),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1", "us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   func() *filterExpr { expr, _ := filter.Parse(`none == ""`); return expr }(),
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						return nil, errors.New("error")
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error in pagination",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						if params.NextToken == nil {
							return nil, errors.New("error")
						}
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group-3"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-2:000000000000:log-group:test-log-group-3"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(3),
									StoredBytes:     aws.Int64(0),
								},
							},
							NextToken: nil,
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cancel",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(ctx context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						<-ctx.Done()
						return nil, ctx.Err()
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateZero,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
			},
			want:    nil,
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
			got, err := man.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.entries != nil {
				SortEntries(got)
			}
			opt := cmp.AllowUnexported(ListEntryData{}, ListEntry{}, entry{})
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("Manager.List() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
