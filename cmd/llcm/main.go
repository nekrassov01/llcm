package main

import (
	"context"
	"os"
)

func main() {
	ctx := context.Background()
	app := newApp(os.Stdout, os.Stderr)
	if err := app.Run(ctx, os.Args); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
