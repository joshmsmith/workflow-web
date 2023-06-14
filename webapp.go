package main

import (
	"log"

	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"webapp/handlers"
	"webapp/transferclient"
)

var CheckTransferTaskQueueTimer = os.Getenv("CHECK_TRANSFER_TASKQUEUE_TIMER")

// Main
func main() {

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", handlers.Home)
	router.HandleFunc("/home", handlers.Home)

	// accounts database handlers
	router.HandleFunc("/accounts", handlers.ListAccounts)
	router.HandleFunc("/showaccount", handlers.ShowAccount)
	router.HandleFunc("/newaccount", handlers.NewAccount)
	router.HandleFunc("/deleteaccount", handlers.DeleteAccount)

	// moneytransfer database handlers
	router.HandleFunc("/transfers", handlers.ListTransfers)
	router.HandleFunc("/showtransfer", handlers.ShowTransfer)
	router.HandleFunc("/newtransfer", handlers.NewTransfer)
	router.HandleFunc("/resettransfer", handlers.ResetTransfer)
	router.HandleFunc("/queryworkflow", handlers.QueryTransferWorkflow)

	// standing order payment handlers
	router.HandleFunc("/listsorders", handlers.ListSOrders)
	router.HandleFunc("/newsorder", handlers.NewSOrder)
	router.HandleFunc("/amendsorder", handlers.AmendSOrder)
	router.HandleFunc("/cancelsorder", handlers.CancelSOrder)

	// Start periodic background transfer table task
	queryDelay, err := strconv.ParseUint(CheckTransferTaskQueueTimer, 20, 64)
	if err != nil {
		queryDelay = 20
	}
	go transferclient.ExecuteCheckTransferTaskCronJob(queryDelay)

	// Serve
	log.Print("Serve Http on 8085")
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8085",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
