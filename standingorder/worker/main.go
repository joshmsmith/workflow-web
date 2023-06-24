package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	so "webapp/standingorder"
  "webapp/utils"
)

func main() {
	log.Printf("%sGo worker starting..%s", so.ColorGreen, so.ColorReset)

	// Load the Temporal Cloud from env
	clientOptions, err := utils.LoadClientOption()
	if err != nil {
		log.Fatalf("Failed to load Temporal Cloud environment: %v", err)
	}

	log.Println("Go worker connecting to server..")

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	log.Println("Go worker initialising..")
	w := worker.New(c, so.StandingOrdersTaskQueueName, worker.Options{})

	log.Println("Go worker registering for Workflow moneytransfer.StandingOrderWorkflow:")
	w.RegisterWorkflow(so.StandingOrderWorkflow)

	// Start listening to the Task Queue.
	log.Printf("%sGo worker listening on %s task queue..%s", so.ColorGreen, so.StandingOrdersTaskQueueName, so.ColorReset)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start StandingOrderWorkflow Worker", err)
	}

	log.Printf("%sGo worker stopped.%s", so.ColorGreen, so.ColorReset)
}
