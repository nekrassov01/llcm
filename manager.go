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
	// Client is the client for CloudWatch Logs.
	*Client `json:"-"`

	// DesiredState is the desired state of the log group.
	DesiredState DesiredState

	// Filters is the expressions to filter the log group entries.
	Filters []Filter

	// Regions is the list of target regions.
	Regions []string

	// desiredState is the native type of DesiredState.
	desiredState *int32

	// filterFunc is the functions to filtering log group entries.
	filterFns []func(*entry) bool

	// sem is the semaphore for concurrent processing.
	sem *semaphore.Weighted

	// ctx is the context for the manager.
	ctx context.Context
}

// NewManager creates a new manager for log group lifecycle management.
func NewManager(ctx context.Context, client *Client) *Manager {
	return &Manager{
		Client:       client,
		DesiredState: -9999,
		Filters:      nil,
		Regions:      DefaultRegions,
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

// String returns the string representation of the log group lifecycle manager.
func (man *Manager) String() string {
	b, _ := json.MarshalIndent(man, "", "  ")
	return string(b)
}
