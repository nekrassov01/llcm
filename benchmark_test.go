package llcm

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"golang.org/x/sync/semaphore"
)

var (
	benchN = 10
	benchR = []string{"us-east-1"}
)

func prepare(n int, regions []string) *Manager {
	var (
		logGroups       = make([]types.LogGroup, n)
		logGroupName    = aws.String("log-group")
		creationTime    = aws.Int64(time.Now().UnixNano() / int64(time.Millisecond))
		retentionInDays = aws.Int32(365)
		storedBytes     = aws.Int64(1024)
		arn             = aws.String("arn:aws:logs:region:account-id:log-group:log-group")
		logGroupClass   = types.LogGroupClassStandard
	)
	return &Manager{
		Client: newMockClient(&mockClient{
			DescribeLogGroupsFunc: func(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
				for i := 0; i < n; i++ {
					logGroups[i] = types.LogGroup{
						LogGroupName:    logGroupName,
						CreationTime:    creationTime,
						RetentionInDays: retentionInDays,
						StoredBytes:     storedBytes,
						Arn:             arn,
						LogGroupClass:   logGroupClass,
					}
				}
				out := &cloudwatchlogs.DescribeLogGroupsOutput{
					LogGroups: logGroups,
				}
				return out, nil
			},
			DeleteLogGroupFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteLogGroupInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error) {
				return &cloudwatchlogs.DeleteLogGroupOutput{}, nil
			},
			PutRetentionPolicyFunc: func(_ context.Context, _ *cloudwatchlogs.PutRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
				return &cloudwatchlogs.PutRetentionPolicyOutput{}, nil
			},
			DeleteRetentionPolicyFunc: func(_ context.Context, _ *cloudwatchlogs.DeleteRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error) {
				return &cloudwatchlogs.DeleteRetentionPolicyOutput{}, nil
			},
		}),
		DesiredState: 365,
		Regions:      regions,
		sem:          semaphore.NewWeighted(NumWorker),
		ctx:          context.Background(),
	}
}

func BenchmarkList(b *testing.B) {
	man := prepare(benchN, benchR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := man.List()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPreview(b *testing.B) {
	man := prepare(benchN, benchR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := man.Preview()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkApply(b *testing.B) {
	man := prepare(benchN, benchR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := man.Apply(io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}
