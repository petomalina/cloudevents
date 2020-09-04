package main

import (
	"context"
	"flag"
	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/flowup/petermalina/services/user/pkg/models"
	"github.com/google/uuid"
	"log"
)

func main() {
	beta := flag.Bool("beta", false, "")
	async := flag.Bool("async", false, "decides if this should be pubsub message")
	project := flag.String("project", "petermalina", "project to set the data to when async=true")
	topic := flag.String("topic", "cloudevents-sink", "pubsub topic to send data to when async=true")
	flag.Parse()

	target := "https://user-ygmoaymzvq-ez.a.run.app"
	if *beta {
		target = "https://beta---user-ygmoaymzvq-ez.a.run.app"
	}
	log.Println("Targetting:", target)

	var err error
	var proto protocol.Sender

	if *async {
		proto, err = cepubsub.New(context.Background(),
			cepubsub.WithProjectID(*project),
			cepubsub.WithTopicID(*topic))
	} else {
		proto, err = cloudevents.NewHTTP(cloudevents.WithTarget(target))
	}
	if err != nil {
		log.Fatalf("Failed to create protocol, %v", err)
	}

	c, err := cloudevents.NewClient(proto,
		cloudevents.WithTimeNow(),
	)

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("com.petomalina.sample.sent")
	event.SetSource("https://github.com/cloudevents/sdk-go/v2/samples/requester")
	err = event.SetData(cloudevents.ApplicationJSON, &models.User{
		Name: "Peter",
	})
	if err != nil {
		panic(err)
	}

	var res cloudevents.Result
	var resEvent *cloudevents.Event
	if *async {
		res = c.Send(context.Background(), event)
	} else {
		resEvent, res = c.Request(context.Background(), event)
		log.Println(resEvent)
	}

	if res != nil && !cloudevents.IsACK(res) {
		panic(res)
	}

	log.Println("Done")
}
