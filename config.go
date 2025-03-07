package llcm

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	// DefaultRegion is the region speficied by default.
	DefaultRegion = "us-east-1"

	// DefaultRegions is the default target regions.
	DefaultRegions = []string{
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}
)

var (
	// nowFunc is used to mock time in tests.
	nowFunc = time.Now

	// allowedRegions is the list of allowed regions.
	allowedRegions = map[string]struct{}{
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
