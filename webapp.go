package main

import (
	"log"

	"net/http"
	"time"

	"github.com/gorilla/mux"

	"webapp/handlers"
	"webapp/transferclient"
)

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

	// Start periodic background transfer table task
	go transferclient.ExecuteCheckTransferTaskCronJob(30)

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
