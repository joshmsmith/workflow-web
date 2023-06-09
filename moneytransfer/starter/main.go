package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "math/rand"
  "time"

  "go.temporal.io/sdk/client"

  mt "webapp/moneytransfer"
)


/* Main */
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
  workflowID := fmt.Sprintf("go-moneytxfr-wkfl-%s", genRandString(5))

  workflowOptions := client.StartWorkflowOptions {
                       ID:        workflowID,
                       TaskQueue: mt.MoneyTransferTaskQueueName,
  }

  // Sample workflow data
  pmnt := &mt.PaymentDetails {
            SourceAccount: "harry",
            TargetAccount: "sally",
            ReferenceID:   "from Go Starter",
            Amount:        100,
  }
  var delay int = 5 // delay between withdraw and deposit for demo purposes (seconds)

  // ExecuteWorkflow mt.Transfer
  log.Println("Starting moneytransfer workflow on", mt.MoneyTransferTaskQueueName, "task queue")

  we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, mt.Transfer, *pmnt, delay)

  if err != nil {
    log.Fatalln("Unable to execute workflow", err)
  }
  log.Printf("%sWorkflow started:%s (WorkflowID: %s, RunID: %s)", mt.ColorYellow, mt.ColorReset, we.GetID(), we.GetRunID())

  // Check workflow status
  var result string

  err = we.Get(context.Background(), &result)

  if err != nil {
    log.Fatalln("Unable get workflow result", err)
  }

  data, err := json.MarshalIndent(result, "", "  ")
  if err != nil {
    log.Fatalln("Unable to format result in JSON format", err)
  }
  log.Printf("%sWorkflow result:%s %s", mt.ColorYellow, mt.ColorReset, string(data))

  // done
  log.Print("Start workflow client done.");
}


/* genRandString */
func genRandString (length int) string {
  rand.Seed(time.Now().UnixNano())
  b := make([]byte, length+2)
  rand.Read(b)
  return fmt.Sprintf("%x", b)[2 : length+2]
}

