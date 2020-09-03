package main

import (
	"context"
	"flag"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/flowup/petermalina/services/user/pkg/models"
	"github.com/google/uuid"
	"log"
)

func main() {
	beta := flag.Bool("beta", false, "")
	flag.Parse()

	target := "https://user-ygmoaymzvq-ez.a.run.app"
	if *beta {
		target = "https://beta---user-ygmoaymzvq-ez.a.run.app"
	}

	log.Println("Targetting:", target)

	p, err := cloudevents.NewHTTP(cloudevents.WithTarget(target))
	if err != nil {
		log.Fatalf("Failed to create protocol, %v", err)
	}

	c, err := cloudevents.NewClient(p,
		cloudevents.WithTimeNow(),
	)

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("com.petomalina.sample.sent")
	event.SetSource("https://github.com/cloudevents/sdk-go/v2/samples/requester")
	err = event.SetData(cloudevents.ApplicationJSON, &models.User{
		Name: "Peto",
	})
	if err != nil {
		panic(err)
	}

	e, res := c.Request(context.Background(), event)
	if !cloudevents.IsACK(res) {
		panic(err)
	}

	log.Println(e)
}
