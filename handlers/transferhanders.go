package handlers

import (
  "log"
  "fmt"
  "strconv"
  "strings"

  "net/http"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "webapp/utils"
  mt "webapp/moneytransfer"
)

type Transfer struct {
  Id          int
  Origin      string
  Destination string
  Amount      float64
  Reference   string
  Status      string
  TWorkflowId string
  TRunId      string
  TTaskQueue  string
  TInfo       string
}


/* ListTransfers */
func ListTransfers(w http.ResponseWriter, r *http.Request) {

  log.Println("ListTransfers: called")

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := `SELECT id,origin,destination,amount,reference,status FROM moneytransfer.transfer order by id desc`
  rows, dberr := dbc.Query(sqlStatement)
  if dberr != nil {
    if dberr == sql.ErrNoRows {
      log.Println("ListTransfers: no entres found")
    } else {
      log.Fatal(dberr)
    }
  }
  defer rows.Close()

  tf := Transfer{}
  transfers := []Transfer{}

  for rows.Next() {
    rows.Scan(&tf.Id, &tf.Origin, &tf.Destination, &tf.Amount, &tf.Reference, &tf.Status)
    transfers = append(transfers, tf)
  }
  //log.Println("ListTransfers: Transfers:", transfers)

  utils.Render(w, "templates/ListTransfers.html", transfers)
}


/* ShowTransfer */
func ShowTransfer (w http.ResponseWriter, r *http.Request) {

  log.Println("ShowTransfer: called")

  // URL Parameters
  var idstr string
  params := r.URL.Query()
  for k, v := range params {
	if k == "id" { idstr = strings.Join(v,"") }
	log.Println("ShowTransfer: url params:", k, " => ", v)
  }
  id, _ := strconv.Atoi(idstr)

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := fmt.Sprintf("select id,origin,destination,amount,reference,status,t_wkfl_id,t_run_id,t_taskqueue,t_info from moneytransfer.transfer where id=%d", id)
  rows, dberr := dbc.Query(sqlStatement)
  if dberr != nil {
    if dberr == sql.ErrNoRows {
      log.Println("ShowTransfer: no transfer entry found")
    } else {
      log.Fatal(dberr)
    }
  }
  defer rows.Close()

  // read entry
  txfr := Transfer{}

  for rows.Next() {
    rows.Scan(&txfr.Id, &txfr.Origin, &txfr.Destination, &txfr.Amount, &txfr.Reference, &txfr.Status, &txfr.TWorkflowId, &txfr.TRunId, &txfr.TTaskQueue, &txfr.TInfo)
  }

  // Display details for requested entry
  utils.Render(w, "templates/ShowTransfer.html", txfr)
}


/* NewTransfer */
func NewTransfer (w http.ResponseWriter, r *http.Request) {

  log.Println("NewTransfer: called")

  log.Println("NewTransfer: method:", r.Method) //get request method
  if r.Method == "GET" {
    utils.Render(w, "templates/NewTransfer.html", nil)
    return
  }

  r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
  amt, _ := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64)
  txfr := Transfer{
    Origin:      r.FormValue("origin"),
    Destination: r.FormValue("destination"),
    Amount:      amt,
    Reference:   r.FormValue("reference"),
  }
  log.Println("NewTransfer: New Transfer Submitted:", txfr)

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := fmt.Sprintf("insert into moneytransfer.transfer (origin, destination, amount, reference, status) values ('%s','%s',%f,'%s','REQUESTED')", txfr.Origin, txfr.Destination, txfr.Amount, txfr.Reference)

  stmtIns, dberr := dbc.Prepare(sqlStatement)
  if dberr != nil {
      log.Fatal("NewTransfer: transfer insert Prepare failed! ", dberr)
  }
  _, dberr = stmtIns.Exec()
  if dberr != nil {
      log.Fatal("NewTransfer: transfer insert Exec failed! ", dberr)
  }
  log.Println("NewTransfer: New Transfer request added to database.")

  // Render acknowledgement page
  utils.Render(w, "templates/NewTransfer.html", struct{ Success bool }{true})
}


/* ResetTransfer */
func ResetTransfer (w http.ResponseWriter, r *http.Request) {

  log.Println("ResetTransfer: called")

  // URL Parameters
  var idstr string
  params := r.URL.Query()
  for k, v := range params {
	if k == "id" { idstr = strings.Join(v,"") }
	log.Println("ResetTransfer: url params:", k, " => ", v)
  }
  id, _ := strconv.Atoi(idstr)

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := fmt.Sprintf("update moneytransfer.transfer set status='REQUESTED', t_wkfl_id='', t_run_id='', t_taskqueue='', t_info='' where id=%d", id)

  stmtIns, dberr := dbc.Prepare(sqlStatement)
  if dberr != nil {
      log.Fatal("ResetTransfer: transfer update Prepare failed! ", dberr)
  }
  _, dberr = stmtIns.Exec()
  if dberr != nil {
      log.Fatal("ResetTransfer: transfer update Exec failed! ", dberr)
  }
  log.Println("ResetTransfer: transfer reset in database.")


  // Render acknowledgement page
  utils.Render(w, "templates/ResetTransfer.html", struct{ Success bool }{true})
}


/* QueryTransferWorkflow */
func QueryTransferWorkflow (w http.ResponseWriter, r *http.Request) {

  log.Println("QueryTransfer: called")

  // URL Parameters
  var idstr string
  params := r.URL.Query()
  for k, v := range params {
	if k == "id" { idstr = strings.Join(v,"") }
	log.Println("ShowTransfer: url params:", k, " => ", v)
  }
  id, _ := strconv.Atoi(idstr)

  // Get database connection
  dbc, _ := utils.GetDBConnection()
  defer dbc.Close()

  sqlStatement := fmt.Sprintf("select id,t_wkfl_id,t_run_id,t_taskqueue from moneytransfer.transfer where id=%d", id)
  rows, dberr := dbc.Query(sqlStatement)
  if dberr != nil {
    if dberr == sql.ErrNoRows {
      log.Println("ShowTransfer: no transfer entry found")
    } else {
      log.Fatal(dberr)
    }
  }
  defer rows.Close()

  // read entry
  wfinfo := mt.WorkflowInfo{}

  for rows.Next() {
    rows.Scan(&wfinfo.Id, &wfinfo.WorkflowID, &wfinfo.RunID, &wfinfo.TaskQueue)
  }

  // Query the workflow
  mt.QueryMoneyTransfer(w, &wfinfo)
}


//Transfer table:
//+-------------+--------------+------+-----+-------------------+-------------------+
//| Field       | Type         | Null | Key | Default           | Extra             |
//+-------------+--------------+------+-----+-------------------+-------------------+
//| id          | int unsigned | NO   | PRI | NULL              | auto_increment    |
//| origin      | varchar(30)  | NO   | MUL | NULL              |                   |
//| destination | varchar(30)  | NO   |     | NULL              |                   |
//| amount      | float        | NO   |     | NULL              |                   |
//| reference   | varchar(30)  | NO   |     | NULL              |                   |
//| status      | varchar(30)  | NO   |     | NULL              |                   |
//| t_wkfl_id   | varchar(50)  | YES  |     | NULL              |                   |
//| t_run_id    | varchar(50)  | YES  |     | NULL              |                   |
//| t_taskqueue | varchar(50)  | YES  |     | NULL              |                   |
//| t_info      | varchar(250) | YES  |     | NULL              |                   |
//| datestamp   | timestamp    | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
//+-------------+--------------+------+-----+-------------------+-------------------+

