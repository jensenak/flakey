package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/jensenak/flakey"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "flakey", worker.Options{})

	w.RegisterWorkflow(flakey.Workflow)
	w.RegisterActivity(flakey.Start)
	w.RegisterActivity(flakey.GetSteps)
	w.RegisterActivity(flakey.RunSteps)
	w.RegisterActivity(flakey.Submit)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
