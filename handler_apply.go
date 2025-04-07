package llcm

import (
	"fmt"
	"io"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// Apply applies the desired state to the log groups.
func (man *Manager) Apply(w io.Writer) (int32, error) {
	var n int32
	err := man.handle(func(man *Manager, entry *entry) error {
		switch man.desiredState {
		case DesiredStateNone:
			return fmt.Errorf("invalid desired state: %q", man.desiredState)
		case DesiredStateZero:
			if err := man.deleteLogGroup(entry.name, entry.Region); err != nil {
				return err
			}
			fmt.Fprintf(w, "deleted log group: %s\n", entry.LogGroupName)
		case DesiredStateInfinite:
			if err := man.deleteRetentionPolicy(entry.name, entry.Region); err != nil {
				return err
			}
			fmt.Fprintf(w, "deleted retention policy: %s\n", entry.LogGroupName)
		default:
			if err := man.putRetentionPolicy(entry.name, entry.Region); err != nil {
				return err
			}
			fmt.Fprintf(w, "updated retention policy: %s\n", entry.LogGroupName)
		}
		atomic.AddInt32(&n, 1)
		return nil
	})
	return n, err
}

// deleteLogGroup deletes the log group.
func (man *Manager) deleteLogGroup(name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: name,
	}
	_, err := man.client.DeleteLogGroup(man.ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}

// deleteRetentionPolicy deletes the retention policy.
func (man *Manager) deleteRetentionPolicy(name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.DeleteRetentionPolicyInput{
		LogGroupName: name,
	}
	_, err := man.client.DeleteRetentionPolicy(man.ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}

// putRetentionPolicy puts the retention policy.
func (man *Manager) putRetentionPolicy(name *string, region string) error {
	opt := func(o *cloudwatchlogs.Options) {
		o.Region = region
		o.Retryer = retryer
	}
	in := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    name,
		RetentionInDays: man.desiredStateNative,
	}
	_, err := man.client.PutRetentionPolicy(man.ctx, in, opt)
	if err != nil {
		return err
	}
	return nil
}
