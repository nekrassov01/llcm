package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nekrassov01/llcm"
)

var client *llcm.Client

func init() {
	cfg, err := llcm.LoadConfig(context.TODO(), "")
	if err != nil {
		log.Fatal(err)
	}
	client = llcm.NewClient(cfg)
	log.Println("client created in init")
}

func handleRequest(ctx context.Context) error {
	log.Println("handleRequest started")
	w := os.Stdout

	// parse desired state passed as string
	d := os.Getenv("DESIRED_STATE")
	desired, err := llcm.ParseDesiredState(d)
	if err != nil {
		return err
	}

	// parse filters passed as comma-separated string
	f := os.Getenv("FILTERS")
	filter, err := llcm.EvaluateFilter(strings.Split(f, ","))
	if err != nil {
		return err
	}

	// initialize the manager
	man := llcm.NewManager(ctx, client)

	// set filter to the manager
	if err := man.SetFilter(filter); err != nil {
		return err
	}

	// set desired state to the manager
	if err := man.SetDesiredState(desired); err != nil {
		return err
	}

	// run apply operation
	n, err := man.Apply(w)
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
