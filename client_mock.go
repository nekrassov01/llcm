package llcm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

var _ API = (*mockClient)(nil)

// mockClient represents a mock client for CloudWatch Logs.
type mockClient struct {
	DescribeLogGroupsFunc     func(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	PutRetentionPolicyFunc    func(ctx context.Context, params *cloudwatchlogs.PutRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error)
	DeleteRetentionPolicyFunc func(ctx context.Context, params *cloudwatchlogs.DeleteRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error)
	DeleteLogGroupFunc        func(ctx context.Context, params *cloudwatchlogs.DeleteLogGroupInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error)
}

// DescribeLogGroups describes the specified log groups.
func (m *mockClient) DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	return m.DescribeLogGroupsFunc(ctx, params, optFns...)
}

// PutRetentionPolicy puts the retention policy of the specified log group.
func (m *mockClient) PutRetentionPolicy(ctx context.Context, params *cloudwatchlogs.PutRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
	return m.PutRetentionPolicyFunc(ctx, params, optFns...)
}

// DeleteRetentionPolicy deletes the specified retention policy.
func (m *mockClient) DeleteRetentionPolicy(ctx context.Context, params *cloudwatchlogs.DeleteRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error) {
	return m.DeleteRetentionPolicyFunc(ctx, params, optFns...)
}

// DeleteLogGroup deletes the specified log group.
func (m *mockClient) DeleteLogGroup(ctx context.Context, params *cloudwatchlogs.DeleteLogGroupInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error) {
	return m.DeleteLogGroupFunc(ctx, params, optFns...)
}

// newMockClient creates a new mock client.
func newMockClient(api API) *Client {
	return &Client{
		api,
	}
}
