package llcm

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var _ aws.RetryerV2 = (*Retryer)(nil)

var (
	retryer       = NewRetryer(retryableFunc, DelayTimeSec)
	retryableFunc = func(err error) bool {
		return strings.Contains(err.Error(), "api error ThrottlingException")
	}
)

// Retryer represents a retryer for avoiding api rate limit exceeded.
type Retryer struct {
	isErrorRetryableFunc func(error) bool
	delayTimeSec         int
}

// NewRetryer creates a new retryer.
func NewRetryer(isErrorRetryableFunc func(error) bool, delayTimeSec int) *Retryer {
	return &Retryer{
		isErrorRetryableFunc: isErrorRetryableFunc,
		delayTimeSec:         delayTimeSec,
	}
}

// IsErrorRetryable checks if the error is retryable.
func (r *Retryer) IsErrorRetryable(err error) bool {
	return r.isErrorRetryableFunc(err)
}

// MaxAttempts returns the maximum number of retry attempts.
func (r *Retryer) MaxAttempts() int {
	return MaxRetryAttempts
}

// RetryDelay returns the delay time for retry.
func (r *Retryer) RetryDelay(int, error) (time.Duration, error) {
	if r.delayTimeSec <= 0 {
		return 0, fmt.Errorf("invalid delay time: %d", r.delayTimeSec)
	}
	var (
		rng  = rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
		wait = 1
	)
	if r.delayTimeSec > 1 {
		wait += rng.Intn(r.delayTimeSec)
	}
	return time.Duration(wait) * time.Second, nil
}

// GetRetryToken returns the retry token.
func (r *Retryer) GetRetryToken(context.Context, error) (func(error) error, error) {
	return func(error) error { return nil }, nil
}

// GetInitialToken returns the initial token.
func (r *Retryer) GetInitialToken() func(error) error {
	return func(error) error { return nil }
}

// GetAttemptToken returns the attempt token.
func (r *Retryer) GetAttemptToken(context.Context) (func(error) error, error) {
	return func(error) error { return nil }, nil
}
