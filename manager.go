package llcm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/sync/semaphore"
)

// NumWorker is the number of workers for concurrent processing.
var NumWorker = int64(runtime.NumCPU()*2 + 1)

// Manager represents a log group lifecycle manager.
type Manager struct {
	*Client            `json:"-"`          // The client for CloudWatch Logs.
	regions            []string            // The list of target regions.
	desiredState       DesiredState        // The desired state of the log group.
	desiredStateNative *int32              // The desired state with the native type.
	filters            []Filter            // The expressions for filtering log groups.
	filterFns          []func(*entry) bool // The list of functions for filtering log groups.
	sem                *semaphore.Weighted // The weighted semaphore for concurrent processing.
	ctx                context.Context     // The context for concurrent processing.
}

// NewManager creates a new manager for log group lifecycle management.
func NewManager(ctx context.Context, client *Client) *Manager {
	return &Manager{
		Client:       client,
		regions:      DefaultRegions,
		desiredState: DesiredStateNone,
		filters:      nil,
		sem:          semaphore.NewWeighted(NumWorker),
		ctx:          ctx,
	}
}

// SetRegion sets the specified regions.
func (man *Manager) SetRegion(regions []string) error {
	if len(regions) == 0 {
		return nil
	}
	for _, region := range regions {
		if _, ok := allowedRegions[region]; !ok {
			return fmt.Errorf("unsupported region: %s", region)
		}
	}
	man.regions = regions
	return nil
}

// SetDesiredState sets the desired state.
func (man *Manager) SetDesiredState(desired DesiredState) error {
	if desired == DesiredStateNone {
		return errors.New("invalid desired state")
	}
	man.desiredState = desired
	man.desiredStateNative = aws.Int32(int32(man.desiredState))
	return nil
}

// SetFilter sets the filter expressions.
func (man *Manager) SetFilter(filters []Filter) error {
	if len(filters) == 0 {
		return nil
	}
	if err := man.setFilter(filters); err != nil {
		return err
	}
	return nil
}

// String returns the string representation of the manager.
func (man *Manager) String() string {
	s := struct {
		Regions      []string `json:"regions"`
		DesiredState string   `json:"desiredState"`
		Filters      []Filter `json:"filters"`
	}{
		Regions:      man.regions,
		DesiredState: man.desiredState.String(),
		Filters:      man.filters,
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}
