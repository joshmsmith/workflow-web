package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	sw "webapp/scheduleworkflow"
  "webapp/utils"
)

func main() {
	log.Printf("%sGo worker starting..%s", sw.ColorGreen, sw.ColorReset)

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
	taskqueuename := "scheduleworkflow"
	w := worker.New(c, taskqueuename, worker.Options{})

	log.Println("Go worker registering for Workflow scheduleworkflow..")
	w.RegisterWorkflow(sw.ScheduleWorkflow)
	log.Println("Go worker registering for Activity ScheduleEmail..")
	w.RegisterActivity(sw.ScheduleEmail)

	log.Printf("%sGo worker listening on %s task queue..%s", sw.ColorGreen, taskqueuename, sw.ColorReset)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	log.Printf("%sGo worker stopped.%s", sw.ColorGreen, sw.ColorReset)
}
