package llcm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/sync/semaphore"
)

// Manager represents a log group lifecycle manager.
type Manager struct {
	*Client      `json:"-"`          // The client for CloudWatch Logs.
	Regions      []string            // The list of target regions.
	DesiredState DesiredState        // The desired state of the log group.
	Filters      []Filter            // The expressions for filtering log groups.
	desiredState *int32              // The desired state with the native type.
	filterFns    []func(*entry) bool // The list of functions for filtering log groups.
	sem          *semaphore.Weighted // The weighted semaphore for concurrent processing.
	ctx          context.Context     // The context for concurrent processing.
}

// NewManager creates a new manager for log group lifecycle management.
func NewManager(ctx context.Context, client *Client) *Manager {
	return &Manager{
		Client:       client,
		Regions:      DefaultRegions,
		DesiredState: -9999,
		Filters:      nil,
		sem:          semaphore.NewWeighted(NumWorker),
		ctx:          ctx,
	}
}

// SetRegions sets the specified regions.
func (man *Manager) SetRegions(regions []string) error {
	if len(regions) == 0 {
		return nil
	}
	for _, region := range regions {
		if !slices.Contains(DefaultRegions, region) {
			return fmt.Errorf("unsupported region: %s", region)
		}
	}
	man.Regions = regions
	return nil
}

// SetDesiredState sets the desired state.
func (man *Manager) SetDesiredState(desired DesiredState) error {
	if desired == DesiredStateNone {
		return errors.New("invalid desired state")
	}
	man.DesiredState = desired
	man.desiredState = aws.Int32(int32(man.DesiredState))
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
	b, _ := json.MarshalIndent(man, "", "  ")
	return string(b)
}
