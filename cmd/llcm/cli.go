package main

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dustin/go-humanize"
	"github.com/nekrassov01/llcm"
	"github.com/nekrassov01/logwrapper/log"
	"github.com/urfave/cli/v2"
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

type app struct {
	*cli.App
	profile  *cli.StringFlag
	loglevel *cli.StringFlag
	region   *cli.StringSliceFlag
	filter   *cli.StringSliceFlag
	desired  *cli.StringFlag
	output   *cli.StringFlag
}

func newApp(w, ew io.Writer) *app {
	logger = log.NewAppLogger(ew, defaultLogLevel, defaultLogStyle, label)
	a := app{}
	a.profile = &cli.StringFlag{
		Name:    "profile",
		Aliases: []string{"p"},
		Usage:   "set aws profile",
		EnvVars: []string{"AWS_PROFILE"},
	}
	a.loglevel = &cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Usage:   "set log level",
		EnvVars: []string{label + "_LOG_LEVEL"},
		Value:   log.InfoLevel.String(),
	}
	a.region = &cli.StringSliceFlag{
		Name:        "region",
		Aliases:     []string{"r"},
		Usage:       "set target regions",
		Value:       cli.NewStringSlice(llcm.DefaultRegions...),
		DefaultText: "all regions with no opt-in required",
	}
	a.filter = &cli.StringSliceFlag{
		Name:    "filter",
		Aliases: []string{"f"},
		Usage:   "set expressions to filter log groups",
	}
	a.desired = &cli.StringFlag{
		Name:     "desired",
		Aliases:  []string{"d"},
		Usage:    "set the desired state",
		Required: true,
	}
	a.output = &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "set output type",
		EnvVars: []string{label + "_OUTPUT_TYPE"},
		Value:   llcm.OutputTypeCompressedText.String(),
	}
	a.App = &cli.App{
		Name:                 name,
		Version:              version(),
		Usage:                "AWS Log groups lifecycle manager",
		Description:          "A listing, updating, and deleting tool to manage the lifecycle of Amazon CloudWatch Logs.\nIt handles multiple regions fast while avoiding throttling. It can also return simulation\nresults based on the desired state.",
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Writer:               w,
		ErrWriter:            ew,
		Metadata:             map[string]any{},
		Commands: []*cli.Command{
			{
				Name:        "completion",
				Usage:       "Generate shell completion script",
				Description: "Generate shell completion script for the specified shell.",
				Action:      a.completion,
			},
			{
				Name:        "list",
				Usage:       "List log group entries with specified format",
				Description: "List collects basic information about log groups from multiple specified regions and\nreturns it in a specified format. ",
				Before:      a.before,
				Action:      a.list,
				Flags:       []cli.Flag{a.profile, a.loglevel, a.region, a.filter, a.output},
			},
			{
				Name:        "preview",
				Usage:       "Preview simulation results based on desired state",
				Description: "Preview performs a simple calculation based on `DesiredState` specified in the argument\nand returns a simulated list including `ReducibleBytes`, `RemainingBytes`, etc.",
				Before:      a.before,
				Action:      a.preview,
				Flags:       []cli.Flag{a.profile, a.loglevel, a.region, a.filter, a.desired, a.output},
			},
			{
				Name:        "apply",
				Usage:       "Apply desired state to log group entries",
				Description: "Apply deletes and updates target log groups in batches based on `DesiredState`.\nIt is fast across multiple regions, but cleverly avoids throttling.",
				Before:      a.before,
				Action:      a.apply,
				Flags:       []cli.Flag{a.profile, a.loglevel, a.region, a.filter, a.desired},
			},
		},
	}
	return &a
}

func (a *app) before(c *cli.Context) error {
	// parses log level passed as string
	level, err := log.ParseLevel(c.String(a.loglevel.Name))
	if err != nil {
		return err
	}

	// set logger for the application
	logger.SetLevel(level)

	// load aws config with the specified profile
	cfg, err := llcm.LoadConfig(c.Context, c.String(a.profile.Name))
	if err != nil {
		return err
	}

	// set logger for the AWS SDK
	cfg.Logger = log.NewSDKLogger(a.ErrWriter, level, defaultLogStyle, "SDK")
	cfg.ClientLogMode = aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogSigning | aws.LogDeprecatedUsage

	// set aws config to the metadata
	a.Metadata["config"] = cfg

	return nil
}

func (a *app) completion(c *cli.Context) error {
	// parse shell passed as string
	n, err := parseShell(c.Args().First())
	if err != nil {
		return err
	}

	// return the completion script
	switch n {
	case bash:
		fmt.Fprintln(a.Writer, bashScript)
	case zsh:
		fmt.Fprintln(a.Writer, zshScript)
	case pwsh:
		fmt.Fprintln(a.Writer, pwshScript)
	default:
	}

	return nil
}

func (a *app) list(c *cli.Context) error {
	// parse output type passed as string
	outputType, err := llcm.ParseOutputType(c.String(a.output.Name))
	if err != nil {
		return err
	}

	// logging at process start
	logger.Info(
		"started",
		"at", time.Now().Format(time.RFC3339),
		"output", outputType,
	)

	// evaluate filter expressions passed as string
	filter, err := llcm.EvaluateFilter(c.StringSlice(a.filter.Name))
	if err != nil {
		return err
	}

	// get aws config from the metadata
	cfg := a.Metadata["config"].(aws.Config)

	// create a new client
	client := llcm.NewClient(cfg)

	// initialize the manager
	man := llcm.NewManager(c.Context, client)

	// set regions to the manager
	if err := man.SetRegions(c.StringSlice(a.region.Name)); err != nil {
		return err
	}

	// set filter to the manager
	if err := man.SetFilter(filter); err != nil {
		return err
	}

	// run list operation
	data, err := man.List()
	if err != nil {
		return err
	}
	debug(man)

	// sort result
	llcm.SortEntries(data.Entries())

	// render result
	ren := llcm.NewRenderer(a.Writer, data, outputType)
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

func (a *app) preview(c *cli.Context) error {
	// parse output type passed as string
	outputType, err := llcm.ParseOutputType(c.String(a.output.Name))
	if err != nil {
		return err
	}

	// logging at process start
	logger.Info(
		"started",
		"at", time.Now().Format(time.RFC3339),
		"desired", c.String(a.desired.Name),
		"output", outputType,
	)

	// evaluate filter expressions passed as string
	filter, err := llcm.EvaluateFilter(c.StringSlice(a.filter.Name))
	if err != nil {
		return err
	}

	// parse desired state passed as string
	desired, err := llcm.ParseDesiredState(c.String(a.desired.Name))
	if err != nil {
		return err
	}

	// get aws config from the metadata
	cfg := a.Metadata["config"].(aws.Config)

	// create a new client
	client := llcm.NewClient(cfg)

	// initialize the manager
	man := llcm.NewManager(c.Context, client)

	// set regions to the manager
	if err := man.SetRegions(c.StringSlice(a.region.Name)); err != nil {
		return err
	}

	// set desired state to the manager
	if err := man.SetDesiredState(desired); err != nil {
		return err
	}

	// set filter to the manager
	if err := man.SetFilter(filter); err != nil {
		return err
	}

	// run preview operation
	data, err := man.Preview()
	if err != nil {
		return err
	}
	debug(man)

	// sort result
	llcm.SortEntries(data.Entries())

	// render result
	ren := llcm.NewRenderer(a.Writer, data, outputType)
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

func (a *app) apply(c *cli.Context) error {
	// logging at process start
	logger.Info(
		"started",
		"at", time.Now().Format(time.RFC3339),
		"desired", c.String(a.desired.Name),
	)

	// evaluate filter expressions passed as string
	filter, err := llcm.EvaluateFilter(c.StringSlice(a.filter.Name))
	if err != nil {
		return err
	}

	// parse desired state passed as string
	desired, err := llcm.ParseDesiredState(c.String(a.desired.Name))
	if err != nil {
		return err
	}

	// get aws config from the metadata
	cfg := a.Metadata["config"].(aws.Config)

	// create a new client
	client := llcm.NewClient(cfg)

	// initialize the manager
	man := llcm.NewManager(c.Context, client)

	// set regions to the manager
	if err := man.SetRegions(c.StringSlice(a.region.Name)); err != nil {
		return err
	}

	// set desired state to the manager
	if err := man.SetDesiredState(desired); err != nil {
		return err
	}

	// set filter to the manager
	if err := man.SetFilter(filter); err != nil {
		return err
	}

	// run apply operation
	n, err := man.Apply(a.Writer)
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

func debug(man *llcm.Manager) {
	logger.Debug("Manager\n" + man.String() + "\n")
}
