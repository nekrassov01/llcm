package main

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dustin/go-humanize"
	"github.com/nekrassov01/llcm"
	"github.com/nekrassov01/logwrapper/log"
	"github.com/urfave/cli/v3"
)

const (
	name  = "llcm"
	label = "LLCM"
)

var (
	logger          = &log.AppLogger{}
	defaultLogLevel = log.InfoLevel
	defaultLogStyle = log.DefaultStyles()
)

func newCmd(w, ew io.Writer) *cli.Command {
	logger = log.NewAppLogger(ew, defaultLogLevel, defaultLogStyle, label)

	profile := &cli.StringFlag{
		Name:    "profile",
		Aliases: []string{"p"},
		Usage:   "set aws profile",
		Sources: cli.EnvVars("AWS_PROFILE"),
	}

	loglevel := &cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Usage:   "set log level",
		Sources: cli.EnvVars(label + "_LOG_LEVEL"),
		Value:   log.InfoLevel.String(),
	}

	region := &cli.StringSliceFlag{
		Name:        "region",
		Aliases:     []string{"r"},
		Usage:       "set target regions",
		Value:       llcm.DefaultRegions,
		DefaultText: "all regions with no opt-in",
	}

	filter := &cli.StringFlag{
		Name:    "filter",
		Aliases: []string{"f"},
		Usage:   "set expressions to filter log groups",
	}

	desired := &cli.StringFlag{
		Name:     "desired",
		Aliases:  []string{"d"},
		Usage:    "set the desired state",
		Required: true,
	}

	output := &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "set output type",
		Sources: cli.EnvVars(label + "_OUTPUT_TYPE"),
		Value:   llcm.OutputTypeCompressedText.String(),
	}

	debug := func(man *llcm.Manager) {
		logger.Debug("Manager\n" + man.String() + "\n")
	}

	newManager := func(cmd *cli.Command) (*llcm.Manager, error) {
		// get aws config from the metadata
		cfg := cmd.Metadata["config"].(aws.Config)

		// create a new client
		client := llcm.NewClient(cfg)

		// initialize the manager
		man := llcm.NewManager(client)

		// set regions to the manager
		if err := man.SetRegion(cmd.StringSlice(region.Name)); err != nil {
			return nil, err
		}

		// set filter to the manager
		if err := man.SetFilter(cmd.String(filter.Name)); err != nil {
			return nil, err
		}

		return man, nil
	}

	before := func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// parse log level passed as string
		level, err := log.ParseLevel(cmd.String(loglevel.Name))
		if err != nil {
			return ctx, err
		}

		// set logger for the application
		logger.SetLevel(level)

		// load aws config with the specified profile
		cfg, err := llcm.LoadConfig(ctx, cmd.String(profile.Name))
		if err != nil {
			return ctx, err
		}

		// set logger for the AWS SDK
		cfg.Logger = log.NewSDKLogger(ew, level, defaultLogStyle, "SDK")
		cfg.ClientLogMode = aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogSigning | aws.LogDeprecatedUsage

		// set aws config to the metadata
		cmd.Metadata["config"] = cfg

		return ctx, nil
	}

	list := func(ctx context.Context, cmd *cli.Command) error {
		// logging at process start
		logger.Info(
			"started",
			"at", time.Now().Format(time.RFC3339),
		)

		// create manager with common settings
		man, err := newManager(cmd)
		if err != nil {
			return err
		}

		// run list operation
		data, err := man.List(ctx)
		if err != nil {
			return err
		}
		debug(man)

		// sort result
		llcm.SortEntries(data)

		// create renderer with data
		ren := llcm.NewRenderer(w, data)

		// set output type passed as string
		if err := ren.SetOutputType(cmd.String(output.Name)); err != nil {
			return err
		}

		// render result
		if err := ren.Render(); err != nil {
			return err
		}

		// logging at process stop with total bytes
		total := data.Total()
		logger.Info(
			"stopped",
			llcm.TotalStoredBytesLabel, humanize.Comma(total[llcm.TotalStoredBytesLabel]),
		)

		return nil
	}

	preview := func(ctx context.Context, cmd *cli.Command) error {
		// logging at process start
		logger.Info(
			"started",
			"at", time.Now().Format(time.RFC3339),
		)

		// create manager with common settings
		man, err := newManager(cmd)
		if err != nil {
			return err
		}

		// set desired state to the manager
		if err := man.SetDesiredState(cmd.String(desired.Name)); err != nil {
			return err
		}

		// run preview operation
		data, err := man.Preview(ctx)
		if err != nil {
			return err
		}
		debug(man)

		// sort result
		llcm.SortEntries(data)

		// create renderer with data
		ren := llcm.NewRenderer(w, data)

		// set output type passed as string
		if err := ren.SetOutputType(cmd.String(output.Name)); err != nil {
			return err
		}

		// render result
		if err := ren.Render(); err != nil {
			return err
		}

		// logging at process stop with the total bytes information
		total := data.Total()
		logger.Info(
			"stopped",
			llcm.TotalStoredBytesLabel, humanize.Comma(total[llcm.TotalStoredBytesLabel]),
			llcm.TotalReducibleBytesLabel, humanize.Comma(total[llcm.TotalReducibleBytesLabel]),
			llcm.TotalRemainingBytesLabel, humanize.Comma(total[llcm.TotalRemainingBytesLabel]),
		)

		return nil
	}

	apply := func(ctx context.Context, cmd *cli.Command) error {
		// logging at process start
		logger.Info(
			"started",
			"at", time.Now().Format(time.RFC3339),
		)

		// create manager with common settings
		man, err := newManager(cmd)
		if err != nil {
			return err
		}

		// set desired state to the manager
		if err := man.SetDesiredState(cmd.String(desired.Name)); err != nil {
			return err
		}

		// run apply operation
		n, err := man.Apply(ctx, w)
		if err != nil {
			return err
		}
		debug(man)

		// logging at process stop with the number of applied entries
		logger.Info(
			"stopped",
			"applied", n,
		)

		return nil
	}

	return &cli.Command{
		Name:                  name,
		Version:               getVersion(),
		Usage:                 "AWS Log groups lifecycle manager",
		Description:           "A listing, updating, and deleting tool to manage the lifecycle of Amazon CloudWatch Logs.\nIt handles multiple regions fast while avoiding throttling. It can also return simulation\nresults based on the desired state.",
		HideHelpCommand:       true,
		EnableShellCompletion: true,
		Writer:                w,
		ErrWriter:             ew,
		Metadata:              map[string]any{},
		Commands: []*cli.Command{
			{
				Name:        "list",
				Usage:       "List log group entries with specified format",
				Description: "List collects basic information about log groups from multiple specified regions and\nreturns it in a specified format.",
				Before:      before,
				Action:      list,
				Flags:       []cli.Flag{profile, loglevel, region, filter, output},
			},
			{
				Name:        "preview",
				Usage:       "Preview simulation results based on desired state",
				Description: "Preview performs a simple calculation based on `DesiredState` specified in the argument\nand returns a simulated list including `ReducibleBytes`, `RemainingBytes`, etc.",
				Before:      before,
				Action:      preview,
				Flags:       []cli.Flag{profile, loglevel, region, filter, desired, output},
			},
			{
				Name:        "apply",
				Usage:       "Apply desired state to log group entries",
				Description: "Apply deletes and updates target log groups in batches based on `DesiredState`.\nIt is fast across multiple regions, but cleverly avoids throttling.",
				Before:      before,
				Action:      apply,
				Flags:       []cli.Flag{profile, loglevel, region, filter, desired},
			},
		},
	}
}
