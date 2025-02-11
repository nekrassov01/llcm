package llcm

import (
	"context"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	// NumWorker is the number of workers for concurrent processing.
	NumWorker = int64(runtime.NumCPU()*2 + 1)

	// MaxRetryAttempts is the maximum number of retry attempts.
	MaxRetryAttempts = 10

	// DelayTimeSec is the sleep time in seconds for retry.
	DelayTimeSec = 3

	// MaxPieChartItems is the maximum number of items in a pie chart.
	MaxPieChartItems = 11

	// MaxBarChartItems is the maximum number of items in a bar chart.
	MaxBarChartItems = 31

	// DefaultRegion is the region speficied by default.
	DefaultRegion = "us-east-1"

	// DefaultRegions is the default target regions.
	DefaultRegions = []string{
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"sa-east-1",
	}
)

var (
	nowFunc             = time.Now
	globalEntriesSize   = 8192
	regionalEntriesSize = 1024
	allowedRegions      = map[string]struct{}{
		"af-south-1":     {},
		"ap-east-1":      {},
		"ap-northeast-1": {},
		"ap-northeast-2": {},
		"ap-northeast-3": {},
		"ap-south-1":     {},
		"ap-south-2":     {},
		"ap-southeast-1": {},
		"ap-southeast-2": {},
		"ap-southeast-3": {},
		"ap-southeast-4": {},
		"ap-southeast-5": {},
		"ap-southeast-7": {},
		"ca-central-1":   {},
		"ca-west-1":      {},
		"eu-central-1":   {},
		"eu-central-2":   {},
		"eu-north-1":     {},
		"eu-south-1":     {},
		"eu-south-2":     {},
		"eu-west-1":      {},
		"eu-west-2":      {},
		"eu-west-3":      {},
		"il-central-1":   {},
		"me-central-1":   {},
		"me-south-1":     {},
		"mx-central-1":   {},
		"sa-east-1":      {},
		"us-east-1":      {},
		"us-east-2":      {},
		"us-west-1":      {},
		"us-west-2":      {},
	}
)

// LoadConfig loads the aws config.
func LoadConfig(ctx context.Context, profile string) (aws.Config, error) {
	var (
		cfg aws.Config
		err error
	)
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return aws.Config{}, err
	}
	if cfg.Region == "" {
		cfg.Region = DefaultRegion
	}
	return cfg, nil
}
