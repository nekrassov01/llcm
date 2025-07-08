package main

import (
	"context"
	"os"
)

func main() {
	ctx := context.Background()
	cmd := newCmd(os.Stdout, os.Stderr)
	if err := cmd.Run(ctx, os.Args); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
