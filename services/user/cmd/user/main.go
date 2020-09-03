package main

import (
	"context"
	"github.com/blendle/zapdriver"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/flowup/petermalina/services/user/pkg/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strings"
)

var (
	L *zap.Logger
)

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

	err = http.ListenAndServe(":"+viper.GetString("port"), httpReceiver)
	if err != nil {
		panic(err)
	}
}

func receive(event cloudevents.Event) *cloudevents.Event {
	L.Info("Received new message", zap.Any("event", event))

	var x models.User
	err := event.DataAs(&x)
	if err != nil {
		L.Fatal(err.Error())
		return nil
	}

	//x.Hash = fmt.Sprintf("%x", sha256.New().Sum(event.Data()))

	err = event.SetData(cloudevents.ApplicationJSON, &x)
	if err != nil {
		L.Fatal(err.Error())
		return nil
	}

	return &event
}
