package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nekrassov01/llcm"
)

var (
	client  *llcm.Client
	filter  string
	desired string
)

func init() {
	f := os.Getenv("FILTER")
	if f == "" {
		log.Fatal("environment variable FILTER not set")
	}
	filter = f

	d := os.Getenv("DESIRED_STATE")
	if d == "" {
		log.Fatal("environment variable DESIRED_STATE not set")
	}
	desired = d

	cfg, err := llcm.LoadConfig(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}

	client = llcm.NewClient(cfg)

	log.Println("client created in init")
}

func handleRequest(ctx context.Context) error {
	log.Println("handleRequest started")
	w := os.Stdout

	// initialize the manager
	man := llcm.NewManager(client)

	// set filter to the manager
	if err := man.SetFilter(filter); err != nil {
		return err
	}

	// set desired state to the manager
	if err := man.SetDesiredState(desired); err != nil {
		return err
	}

	// run apply operation
	n, err := man.Apply(ctx, w)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "done: %d\n", n)

	log.Println("handleRequest finished")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
