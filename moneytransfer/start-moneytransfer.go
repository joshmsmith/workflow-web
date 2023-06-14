package moneytransfer

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "math/rand"

  "go.temporal.io/sdk/client"
)

/*
 * StartMoneyTransfer - App entry point to run Temporal Workflows
 */
func StartMoneyTransfer(pmnt *PaymentDetails) (wfinfo *WorkflowInfo, starterr error) {

  thisid := fmt.Sprint(rand.Intn(99999))
  log.Printf("StartMoneyTransfer-%s: called, PaymentDetails: %#v", thisid, *pmnt)

  // Initialise return object
  wfinfo = &WorkflowInfo{
    Id:         0,
    WorkflowID: fmt.Sprintf("go-txfr-webtask-wkfl-%s", thisid),
    RunID:      "",
    TaskQueue:  MoneyTransferTaskQueueName,
    Info:       "",
    Status:     "ERROR",
  }

  // Load the Temporal Cloud from env
  clientOptions, err := LoadClientOption()
  if err != nil {
    log.Printf("StartMoneyTransfer-%s: Failed to load Temporal Cloud environment: %v", thisid, err)
    wfinfo.Info = err.Error()
    return wfinfo, err
  }
  c, err := client.Dial(clientOptions)
  if err != nil {
    log.Printf("StartMoneyTransfer-%s: Unable to create Temporal client: %v", thisid, err)
    wfinfo.Info = err.Error()
    return wfinfo, err
  }
  defer c.Close()

  // Temporal Client Start Workflow Options
  workflowID := fmt.Sprintf("go-txfr-webtask-wkfl-%s", thisid)

  workflowOptions := client.StartWorkflowOptions{
    ID:        workflowID,
    TaskQueue: MoneyTransferTaskQueueName,
  }

  var delay int = 10 // delay between withdraw and deposit for demo purposes (seconds)

  // ExecuteWorkflow moneytransfer.Transfer
  log.Printf("StartMoneyTransfer-%s: Starting moneytransfer workflow on %s task queue", thisid, MoneyTransferTaskQueueName)

  we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, Transfer, *pmnt, delay)

  if err != nil {
    log.Printf("StartMoneyTransfer-%s: Error, Unable to execute workflow %v", thisid, err)
    wfinfo.Info = err.Error()
    return wfinfo, err
  }
  wfinfo.WorkflowID = we.GetID()
  wfinfo.RunID = we.GetRunID()
  log.Printf("StartMoneyTransfer-%s: %sStarted workflow: WorkflowID: %s, RunID: %s%s",
    thisid, ColorYellow, wfinfo.WorkflowID, wfinfo.RunID, ColorReset)

  // Check workflow status
  var result string

  err = we.Get(context.Background(), &result)

  if err != nil {
    log.Printf("StartMoneyTransfer-%s: %sWorkflow returned failure:%s %v", thisid, ColorRed, ColorReset, err)
    wfinfo.Info = err.Error()
    return wfinfo, err
  }

  data, err := json.MarshalIndent(result, "", "  ")
  if err != nil {
    log.Printf("StartMoneyTransfer-%s: Unable to format result in JSON format: %v", thisid, err)
    wfinfo.Info = err.Error()
    return wfinfo, err
  }
  log.Printf("StartMoneyTransfer-%s: %sWorkflow result: %s%s", thisid, ColorYellow, string(data), ColorReset)
  wfinfo.Info = trimQuotes(string(data))
  wfinfo.Status = "COMPLETED"

  // done
  log.Printf("StartMoneyTransfer-%s: done.", thisid)

  return wfinfo, err
}

/* trimQuotes from string */
func trimQuotes(s string) string {
  if len(s) >= 2 {
    if s[0] == '"' && s[len(s)-1] == '"' {
      return s[1 : len(s)-1]
    }
  }
  return s
}
