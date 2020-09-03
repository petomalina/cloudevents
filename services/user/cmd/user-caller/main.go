package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"log"
)

type User struct {
	Name string
}

func main() {
	p, err := cloudevents.NewHTTP(cloudevents.WithTarget("https://user-ygmoaymzvq-ez.a.run.app"))
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
	err = event.SetData(cloudevents.ApplicationJSON, &User{
		Name: "Peto",
	})
	if err != nil {
		panic(err)
	}

	res := c.Send(context.Background(), event)
	if !cloudevents.IsACK(res) {
		panic(err)
	}
}
