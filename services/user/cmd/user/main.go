package main

import (
	"context"
	"github.com/blendle/zapdriver"
	cloudevents "github.com/cloudevents/sdk-go/v2"
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

type User struct {
	Name string
}

func receive(event cloudevents.Event) cloudevents.Result {
	L.Info("Received new message", zap.Any("event", event))

	var x User
	err := event.DataAs(&x)
	if err != nil {
		return err
	}

	return nil
}
