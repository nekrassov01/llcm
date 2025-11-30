package llcm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

var _ API = (*Client)(nil)

// API represents an interface for CloudWatch Logs.
type API interface {
	DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	PutRetentionPolicy(ctx context.Context, params *cloudwatchlogs.PutRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error)
	DeleteRetentionPolicy(ctx context.Context, params *cloudwatchlogs.DeleteRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error)
	DeleteLogGroup(ctx context.Context, params *cloudwatchlogs.DeleteLogGroupInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogGroupOutput, error)
	PutLogGroupDeletionProtection(ctx context.Context, params *cloudwatchlogs.PutLogGroupDeletionProtectionInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogGroupDeletionProtectionOutput, error)
}

// Client represents a client for CloudWatch Logs.
type Client struct {
	API
}

// NewClient creates a new client.
func NewClient(cfg aws.Config) *Client {
	return &Client{
		cloudwatchlogs.NewFromConfig(cfg),
	}
}
