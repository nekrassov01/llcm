package llcm

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/nekrassov01/filter"
	"golang.org/x/sync/semaphore"
)

// NumWorker is the number of workers for concurrent processing.
var NumWorker = int64(runtime.NumCPU()*2 + 1)

type (
	filterExpr   = *filter.Expr  // filterExpr is a type alias for filter.Expr.
	filterTarget = filter.Target // filterTarget is a type alias for filter.Target.
)

// Manager represents a log group lifecycle manager.
type Manager struct {
	client             *Client             // The client for CloudWatch Logs.
	regions            []string            // The list of target regions.
	desiredState       DesiredState        // The desired state of the log group.
	desiredStateNative *int32              // The desired state with the native type.
	filterExpr         filterExpr          // The expressions for filtering log groups.
	filterRaw          string              // The raw filter string.
	sem                *semaphore.Weighted // The weighted semaphore for concurrent processing.
}

// NewManager creates a new manager for log group lifecycle management.
func NewManager(client *Client) *Manager {
	return &Manager{
		client:       client,
		regions:      DefaultRegions,
		desiredState: DesiredStateNone,
		sem:          semaphore.NewWeighted(NumWorker),
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
func (man *Manager) SetDesiredState(desired string) error {
	d, err := ParseDesiredState(desired)
	if err != nil {
		return err
	}
	man.desiredState = d
	man.desiredStateNative = aws.Int32(int32(man.desiredState))
	return nil
}

// SetFilter sets the filter expressions.
func (man *Manager) SetFilter(raw string) error {
	if raw == "" {
		return nil
	}
	expr, err := filter.Parse(raw)
	if err != nil {
		return fmt.Errorf("failed to parse filter: %w", err)
	}
	man.filterExpr = expr
	man.filterRaw = raw
	return nil
}

// String returns the string representation of the manager.
func (man *Manager) String() string {
	s := struct {
		Regions      []string `json:"regions"`
		DesiredState string   `json:"desiredState"`
		Filter       string   `json:"filter"`
	}{
		Regions:      man.regions,
		DesiredState: man.desiredState.String(),
		Filter:       man.filterRaw,
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}
