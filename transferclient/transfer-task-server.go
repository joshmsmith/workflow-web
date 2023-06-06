package transferclient

import (
  "fmt"
  "log"
  "strings"

  "database/sql"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jasonlvhit/gocron"

  "webapp/utils"

  "webapp/moneytransfer"
)

type Transfer struct {
  Id          int
  Origin      string
  Destination string
  Amount      float64
  Reference   string
  Status      string
}


/* ExecuteCheckTransferTaskCronJob - Cron to call periodic task simulate event queue */
func ExecuteCheckTransferTaskCronJob (internalSeconds uint64) {

  gocron.Every(internalSeconds).Second().Do(CheckTransferQueueTask)
  <-gocron.Start()
}


/* CheckTransferQueueTask - Check transfer queue table task */
func CheckTransferQueueTask() {

  log.Println("CheckTransferQueueTask: called")

  // Call handler to read db task queue and return oldest REQUESTED task
  txfr, err := QueryTransferRequest()
  if err != nil {
    log.Println("CheckTransferQueueTask: Failed to query transfer task queue!", err)
    return
  } else if *txfr == (Transfer{}) {
    // no entry found on queue
    log.Println("CheckTransferQueueTask: No transfers in queue.")
    return
  }

  // transfer to process
  log.Println("CheckTransferQueueTask: Transfer:", *txfr)

  // Populate PaymentDetails object from tranfer task entry
  pmnt := &moneytransfer.PaymentDetails{
    SourceAccount: txfr.Origin,
    TargetAccount: txfr.Destination,
    ReferenceID:   txfr.Reference,
    Amount:        int(txfr.Amount),
  }
  log.Printf("CheckTransferQueueTask: %sPaymentDetails: %v%s", utils.ColorYellow, *pmnt, utils.ColorReset)

  // Call StartMoneyTransfer to start the workflow..
  wfinfo, wferr := moneytransfer.StartMoneyTransfer(pmnt)
  wfinfo.Id = txfr.Id
  if wferr != nil {
    wfinfo.Status = "FAILED"
    if strings.Contains(wferr.Error(), "Insufficient Funds") {
			log.Printf("CheckTransferQueueTask: %sOrigin Account has Insufficient Funds%s", utils.ColorRed, utils.ColorReset)
			wfinfo.Info = "Origin Account has Insufficient Funds"

    } else if strings.Contains(wferr.Error(), "InvalidAccountError") {
      if strings.Contains(wferr.Error(), "Withdraw") {
        log.Printf("CheckTransferQueueTask: %sOrigin Account Invalid%s", utils.ColorRed, utils.ColorReset)
        wfinfo.Info = "Origin Account is Invalid"
      } else if strings.Contains(wferr.Error(), "Deposit") {
        log.Printf("CheckTransferQueueTask: %sDestination Account Invalid, Origin Account Refunded%s", utils.ColorRed, utils.ColorReset)
        wfinfo.Info = "Destination Account is Invalid, Origin Account Refunded"
      } else {
        log.Printf("CheckTransferQueueTask: %sAccount Invalid%s", utils.ColorRed, utils.ColorReset)
        wfinfo.Info = "Invalid Account Error"
      }
    } else {
      log.Printf("CheckTransferQueueTask: %sWorkflow returned error: %v%s", utils.ColorRed, wferr, utils.ColorReset)
    }

    // Update db entry as FAILED
    _ = UpdateTransferRequest(wfinfo)
    log.Printf("CheckTransferQueueTask: %sWorkflow: %s Failed%s", utils.ColorRed, wfinfo.WorkflowID, utils.ColorReset)

  } else {

    // Update db entry as COMPLETED
    wfinfo.Status = "COMPLETED"
    _ = UpdateTransferRequest(wfinfo)
    log.Printf("CheckTransferQueueTask: %sWorkflow: %s Completed%s", utils.ColorYellow, wfinfo.WorkflowID, utils.ColorReset)
  }

  return
}


/* QueryTransferRequest 
 *
 * Query database for a "REQUESTED" task entry
 * Pick oldest entry, update it to "PROCESSING" and return details
 * Only return one entry
 *
 * Multi-sql statement block (dbconnection: ?multiStatements=true)
 *
 * called from: CheckTransferQueueTask
 */
func QueryTransferRequest () (*Transfer, error) {

  log.Println("QueryTransferRequest: called")

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := `set @updatedid:=NULL; update moneytransfer.transfer set status="PROCESSING", id=(@updatedid:=id) where status='REQUESTED' order by id limit 1; select id,origin,destination,amount,reference,status from moneytransfer.transfer where id = @updatedid;`
  rows, dberr := dbc.Query(sqlStatement)
  if dberr != nil {
    if dberr == sql.ErrNoRows {
      log.Println("QueryTransferRequest: no entres found")
      return nil, nil
    } else {
      log.Println("QueryTransferRequest: Query failed!", dberr)
      return nil, dberr
    }
  }

  tf := &Transfer{}

  for rows.Next() {
    rows.Scan(&tf.Id, &tf.Origin, &tf.Destination, &tf.Amount, &tf.Reference, &tf.Status)
  }
  if *tf != (Transfer{}) {
    log.Println("QueryTransferRequest: Transfer:", *tf)
  }

  return tf, nil
}


/* UpdateTransferRequest
 *
 * Workflow has finished so update the database entry with details
 *
 * called from: CheckTransferQueueTask
 */
func UpdateTransferRequest (wfinfo *moneytransfer.WorkflowInfo) error {

  log.Println("UpdateTransferRequest: called (Id:", wfinfo.Id, wfinfo.Status, ")")

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := fmt.Sprintf("update moneytransfer.transfer set status='%s',t_wkfl_id='%s',t_run_id='%s',t_taskqueue='%s',t_info='%s' where id=%d", wfinfo.Status, wfinfo.WorkflowID, wfinfo.RunID, wfinfo.TaskQueue, wfinfo.Info, wfinfo.Id)
  stmtIns, dberr := dbc.Prepare(sqlStatement)
  if dberr != nil {
    log.Println("UpdateTransferRequest: Prepare failed! ", dberr)
    return dberr
  }
  _, dberr = stmtIns.Exec()
  if dberr != nil {
    log.Println("UpdateTransferRequest: update Exec failed! ", dberr)
    return dberr
  }

  return nil
}

