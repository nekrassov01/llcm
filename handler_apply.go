package llcm

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// Apply applies the desired state to the log groups.
func (man *Manager) Apply(ctx context.Context, w io.Writer) (int32, error) {
	var n atomic.Int32
	fn := func(entry *entry) error {
		switch man.desiredState {
		case DesiredStateNone:
			return fmt.Errorf("invalid desired state: %q", man.desiredState)
		case DesiredStateZero:
			if err := man.deleteLogGroup(ctx, entry.name, entry.Region); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(w, "deleted log group: %s\n", entry.LogGroupName)
		case DesiredStateInfinite:
			if err := man.deleteRetentionPolicy(ctx, entry.name, entry.Region); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(w, "deleted retention policy: %s\n", entry.LogGroupName)
		case DesiredStateProtected, DesiredStateUnprotected:
			if err := man.putLogGroupDeletionProtection(ctx, entry.name, entry.Region); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(w, "%s log group: %s\n", man.desiredState.String(), entry.LogGroupName)
		default:
			if err := man.putRetentionPolicy(ctx, entry.name, entry.Region); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(w, "updated retention policy: %s\n", entry.LogGroupName)
		}
		n.Add(1)
		return nil
	}
	err := man.handle(ctx, fn)
	return n.Load(), err
}

// deleteLogGroup deletes the log group.
func (man *Manager) deleteLogGroup(ctx context.Context, name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: name,
	}
	_, err := man.client.DeleteLogGroup(ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}

// deleteRetentionPolicy deletes the retention policy.
func (man *Manager) deleteRetentionPolicy(ctx context.Context, name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.DeleteRetentionPolicyInput{
		LogGroupName: name,
	}
	_, err := man.client.DeleteRetentionPolicy(ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}

// putLogGroupDeletionProtection puts the log group deletion protection.
func (man *Manager) putLogGroupDeletionProtection(ctx context.Context, name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.PutLogGroupDeletionProtectionInput{
		LogGroupIdentifier:        name,
		DeletionProtectionEnabled: man.deletionProtection,
	}
	_, err := man.client.PutLogGroupDeletionProtection(ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}

// putRetentionPolicy puts the retention policy.
func (man *Manager) putRetentionPolicy(ctx context.Context, name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    name,
		RetentionInDays: man.desiredStateNative,
	}
	_, err := man.client.PutRetentionPolicy(ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}
