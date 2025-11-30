package llcm

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/semaphore"
)

func TestManager_Preview(t *testing.T) {
	type fields struct {
		client             *Client
		regions            []string
		desiredState       DesiredState
		desiredStateNative *int32
		filterExpr         filterExpr
		sem                *semaphore.Weighted
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *PreviewEntryData
		wantErr bool
	}{
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
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneMonth,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 600,
				TotalRemainingBytes: 300,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 60,
						ReducibleBytes:  600,
						RemainingBytes:  300,
					},
				},
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
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
								},
								{
									LogGroupName:    aws.String("test-log-group-2"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group-2"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateTwoMonths)),
									StoredBytes:     aws.Int64(1200),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneMonth,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    2100,
				TotalReducibleBytes: 1200,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group-2",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateTwoMonths),
							StoredBytes:     1200,
							name:            aws.String("test-log-group-2"),
						},
						BytesPerDay:     20,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 30,
						ReducibleBytes:  600,
						RemainingBytes:  600,
					},
					{
						entry: &entry{
							LogGroupName:    "test-log-group-1",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group-1"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 60,
						ReducibleBytes:  600,
						RemainingBytes:  300,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "zero bytes",
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
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(0),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneMonth,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    0,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     0,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     0,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "infinite retention",
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
									RetentionInDays: aws.Int32(int32(DesiredStateInfinite)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneMonth,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 600,
				TotalRemainingBytes: 300,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 60,
						ReducibleBytes:  600,
						RemainingBytes:  300,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "infinite retention and zero elapsed days",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-04-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateInfinite)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneMonth,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-04-01T00:00:00Z"),
							ElapsedDays:     0,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     900,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  900,
					},
				},
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
									RetentionInDays: aws.Int32(int32(DesiredStateZero)),
									StoredBytes:     aws.Int64(900),
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
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 900,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 90,
						ReducibleBytes:  900,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "zero retention and zero elapsed days",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-04-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateZero)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateOneMonth,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-04-01T00:00:00Z"),
							ElapsedDays:     0,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     900,
						DesiredState:    int64(DesiredStateOneMonth),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  900,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "desired infinite retention",
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
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateInfinite,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateInfinite),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  900,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "desired infinite retention and zero elapsed days",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-04-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateInfinite,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-04-01T00:00:00Z"),
							ElapsedDays:     0,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     900,
						DesiredState:    int64(DesiredStateInfinite),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  900,
					},
				},
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
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateProtected,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateProtected),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  900,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "deletion protection disabled",
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
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
								},
							},
						}
						return out, nil
					},
				}),
				regions:      []string{"us-east-1"},
				desiredState: DesiredStateUnprotected,
				filterExpr:   nil,
				sem:          semaphore.NewWeighted(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 0,
				TotalRemainingBytes: 900,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateUnprotected),
						ReductionInDays: 0,
						ReducibleBytes:  0,
						RemainingBytes:  900,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "desired zero retention",
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
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
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
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 900,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     10,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 90,
						ReducibleBytes:  900,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "desired zero retention and zero elapsed days",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						out := &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-04-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateThreeMonths)),
									StoredBytes:     aws.Int64(900),
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
			want: &PreviewEntryData{
				TotalStoredBytes:    900,
				TotalReducibleBytes: 900,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-04-01T00:00:00Z"),
							ElapsedDays:     0,
							RetentionInDays: int64(DesiredStateThreeMonths),
							StoredBytes:     900,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     900,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 90,
						ReducibleBytes:  900,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "reducible bytes exceed stored bytes",
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
									RetentionInDays: aws.Int32(int32(DesiredStateOneDay)),
									StoredBytes:     aws.Int64(100),
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
			want: &PreviewEntryData{
				TotalStoredBytes:    100,
				TotalReducibleBytes: 100,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateOneDay),
							StoredBytes:     100,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     100,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 1,
						ReducibleBytes:  100,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "reduction in days convert to 1",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						return &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-04-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateInfinite)),
									StoredBytes:     aws.Int64(100),
								},
							},
						}, nil
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
			want: &PreviewEntryData{
				TotalStoredBytes:    100,
				TotalReducibleBytes: 100,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-04-01T00:00:00Z"),
							ElapsedDays:     0,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     100,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     100,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 1,
						ReducibleBytes:  100,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "bytes per day convert 0 to 1",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						return &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateInfinite)),
									StoredBytes:     aws.Int64(10),
								},
							},
						}, nil
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
			want: &PreviewEntryData{
				TotalStoredBytes:    10,
				TotalReducibleBytes: 10,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     10,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     1,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 90,
						ReducibleBytes:  10,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "bytes per day == stored bytes",
			fields: fields{
				client: newMockClient(&mockClient{
					DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
						return &cloudwatchlogs.DescribeLogGroupsOutput{
							LogGroups: []types.LogGroup{
								{
									LogGroupName:    aws.String("test-log-group"),
									LogGroupArn:     aws.String("arn:aws:logs:us-east-1:123456789012:log-group:test-log-group"),
									LogGroupClass:   types.LogGroupClassStandard,
									CreationTime:    aws.Int64(mustUnixMilli("2025-01-01T00:00:00Z")),
									RetentionInDays: aws.Int32(int32(DesiredStateInfinite)),
									StoredBytes:     aws.Int64(90),
								},
							},
						}, nil
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
			want: &PreviewEntryData{
				TotalStoredBytes:    90,
				TotalReducibleBytes: 90,
				TotalRemainingBytes: 0,
				header:              previewEntryDataHeader,
				entries: []*PreviewEntry{
					{
						entry: &entry{
							LogGroupName:    "test-log-group",
							Region:          "us-east-1",
							Class:           types.LogGroupClassStandard,
							CreatedAt:       mustTime("2025-01-01T00:00:00Z"),
							ElapsedDays:     90,
							RetentionInDays: int64(DesiredStateInfinite),
							StoredBytes:     90,
							name:            aws.String("test-log-group"),
						},
						BytesPerDay:     1,
						DesiredState:    int64(DesiredStateZero),
						ReductionInDays: 90,
						ReducibleBytes:  90,
						RemainingBytes:  0,
					},
				},
			},
			wantErr: false,
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
			got, err := man.Preview(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Preview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.entries != nil {
				SortEntries(got)
			}
			opt := cmp.AllowUnexported(PreviewEntryData{}, PreviewEntry{}, entry{})
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("Manager.Preview() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
