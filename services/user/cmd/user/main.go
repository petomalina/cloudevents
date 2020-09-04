package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/blendle/zapdriver"
	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/flowup/petermalina/services/user/pkg/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"strings"
)

var (
	L *zap.Logger
)

type PushMessage struct {
	Message      pubsub.Message `json:"message"`
	Subscription string         `json:"subscription"`
}

func main() {
	ctx := context.Background()

	viper.SetDefault("log.level", zap.DebugLevel)
	viper.SetDefault("port", 8080)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	config := zapdriver.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.Level(viper.GetInt("log.level")))

	var err error
	L, err = config.Build(zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName("user"),
	))
	if err != nil {
		panic(err)
	}

	httpProto, err := cloudevents.NewHTTP()
	if err != nil {
		panic(err)
	}

	httpReceiver, err := cloudevents.NewHTTPReceiveHandler(ctx, httpProto, receive)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		L.Info("Routing new message with headers:", zap.Any("headers", r.Header))

		if strings.HasPrefix(r.Header.Get("user-agent"), "APIs-Google") {
			msg := PushMessage{}
			err = json.NewDecoder(r.Body).Decode(&msg)
			if err != nil {
				L.Error("An error occurred when unmarhsalling push message:", zap.Error(err))
				return
			}

			ceMsg := cepubsub.NewMessage(&msg.Message)
			ceEvent, err := binding.ToEvent(context.Background(), ceMsg)
			if err != nil {
				L.Error("An error occurred when binding message to event:", zap.Error(err))
				return
			}
			receive(*ceEvent)

			L.Info("PubSubMessage Received")
		} else {
			httpReceiver.ServeHTTP(w, r)
		}
	})

	err = http.ListenAndServe(":"+viper.GetString("port"), handlers.LoggingHandler(os.Stdout, router))
	if err != nil {
		panic(err)
	}
}

func receive(event cloudevents.Event) *cloudevents.Event {
	L.Info("Received new event", zap.Any("event", event))

	var x models.User
	err := event.DataAs(&x)
	if err != nil {
		L.Fatal(err.Error())
		return nil
	}

	x.Hash = fmt.Sprintf("%x", sha256.New().Sum(event.Data()))

	err = event.SetData(cloudevents.ApplicationJSON, &x)
	if err != nil {
		L.Fatal(err.Error())
		return nil
	}

	return &event
}
