package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"go.temporal.io/sdk/client"

	u "webapp/utils"
)

func main() {

	wkflId, err := cancelparseCLIArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("Parameter --workflow-id <workflow id> is required")
	}

	clientOptions, err := u.LoadClientOptions(u.NoSDKMetrics)
	if err != nil {
		log.Fatalf("Failed to load Temporal Cloud environment: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	log.Println("Sending cancelsorder signal to workflow: ", *wkflId)

	// Signal the Workflow Executions to cancel the standing order
	err = c.SignalWorkflow(context.Background(), *wkflId, "", "cancelsorder", true)
	if err != nil {
		log.Fatalln("Unable to signal workflow", err)
	}
}

func cancelparseCLIArgs(args []string) (*string, error) {
	set := flag.NewFlagSet("hello-workflow", flag.ExitOnError)
	wkflId := set.String("workflow-id", "", "Workflow Id to access")
	if err := set.Parse(args); err != nil {
		return nil, fmt.Errorf("failed parsing args: %w", err)
	} else if *wkflId == "" {
		return nil, fmt.Errorf("--workflow-id argument is required")
	}
	return wkflId, nil
}
