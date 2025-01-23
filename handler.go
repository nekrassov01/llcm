package llcm

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// handle enumerates log groups for all regions to get targets for the process.
// For each entry, the specified handler is executed.
func (man *Manager) handle(handler func(*Manager, *entry) error) error {
	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(man.ctx)
		errorChan   = make(chan error, len(man.Regions))
		linked      = aws.Bool(true)
	)
	defer cancel()
	errorFunc := func(err error) {
		select {
		case errorChan <- err:
			cancel()
		default:
		}
	}
	for _, region := range man.Regions {
		region := region
		wg.Add(1)
		go func() {
			defer wg.Done()
			opt := func(o *cloudwatchlogs.Options) {
				o.Region = region
			}
			in := &cloudwatchlogs.DescribeLogGroupsInput{
				NextToken:             nil,
				AccountIdentifiers:    nil,
				IncludeLinkedAccounts: linked,
			}
			for {
				out, err := man.DescribeLogGroups(ctx, in, opt)
				if err != nil {
					errorFunc(err)
					return
				}
				for _, logGroup := range out.LogGroups {
					logGroup := logGroup
					if err := man.sem.Acquire(ctx, 1); err != nil {
						errorFunc(err)
						return
					}
					wg.Add(1)
					go func() {
						defer wg.Done()
						defer man.sem.Release(1)
						entry := newEntry(logGroup, region)
						if !man.applyFilter(entry) {
							return // skip if false
						}
						if err := handler(man, entry); err != nil {
							errorFunc(err)
							return
						}
					}()
				}
				if out.NextToken == nil {
					break
				}
				in.NextToken = out.NextToken
			}
		}()
	}
	go func() {
		wg.Wait()
		close(errorChan)
	}()
	select {
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// applyFilter applies the filter to the entry.
func (man *Manager) applyFilter(entry *entry) bool {
	if len(man.Filters) == 0 {
		return true
	}
	for _, fn := range man.filterFns {
		if !fn(entry) {
			return false
		}
	}
	return true
}

// newEntry creates a new entry from the log group and specified region.
func newEntry(logGroup types.LogGroup, region string) *entry {
	e := &entry{}
	e.LogGroupName = aws.ToString(logGroup.LogGroupName)
	e.Region = region
	e.Source = source(logGroup.LogGroupArn)
	e.Class = logGroup.LogGroupClass
	e.CreatedAt = createdAt(logGroup.CreationTime)
	e.ElapsedDays = elapsedDays(e.CreatedAt)
	e.RetentionInDays = retentionInDays(logGroup.RetentionInDays)
	e.StoredBytes = aws.ToInt64(logGroup.StoredBytes)
	e.name = logGroup.LogGroupName
	return e
}

// source returns the linked source account of the log group.
func source(s *string) string {
	parsed, err := arn.Parse(aws.ToString(s))
	if err != nil {
		return ""
	}
	if parsed.AccountID == "" || parsed.Region == "" {
		return ""
	}
	return parsed.AccountID + "/" + parsed.Region
}

// createdAt returns the creation time of the log group.
func createdAt(t *int64) time.Time {
	return time.Unix(0, aws.ToInt64(t)*int64(time.Millisecond))
}

// elapsedDays returns the elapsed days from the creation time.
func elapsedDays(t time.Time) int64 {
	return int64(nowFunc().Sub(t).Hours() / 24)
}

// retentionInDays returns the retention days from the log group.
func retentionInDays(n *int32) int64 {
	d := aws.ToInt32(n)
	// convert 0 meaning none to 9999 meaning infinite
	if d == 0 {
		d = 9999
	}
	return int64(d)
}
