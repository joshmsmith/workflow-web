package main

import (
  "fmt"
	"context"
	"log"
  "math/rand"
	"time"

	"go.temporal.io/sdk/client"

	mt "webapp/moneytransfer"
)

func main() {

	log.Println("workflow start program..")

	// Load the Temporal Cloud from env
	clientOptions, err := mt.LoadClientOption()
	if err != nil {
		log.Fatalf("Failed to load Temporal Cloud environment: %v", err)
	}
	log.Println("connecting to temporal server..")
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	// Temporal Client Start Workflow Options
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("go-txfr-sorder-wkfl-%d", rand.Intn(99999)),
		TaskQueue: mt.StandingOrdersTaskQueueName,
	}

	// Sample workflow data
	pmnt := &mt.PaymentDetails{
		SourceAccount: "harry",
		TargetAccount: "sally",
		ReferenceID:   "StandingOrder",
		Amount:        11,
	}
	schl := &mt.PaymentSchedule{
		PeriodDuration: time.Duration(30) * time.Second,
		Active:         true,
	}

	log.Println("Starting standingorder workflow ..")
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, mt.StandingOrderWorkflow, *pmnt, *schl)

	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Printf("%sWorkflow started:%s (WorkflowID: %s, RunID: %s)", mt.ColorYellow, mt.ColorReset, we.GetID(), we.GetRunID())
}

