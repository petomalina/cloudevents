package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	firebase "firebase.google.com/go"
	"github.com/blendle/zapdriver"
	v1 "github.com/flowup/petermalina/apis/go-sdk/user/v1"
	"github.com/flowup/petermalina/services/user/internal"
	"github.com/flowup/petermalina/services/user/internal/user"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/reflection"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"strings"
)

func main() {
	ctx := context.Background()
	viper.SetDefault("port", "8080")
	viper.SetDefault("project.id", "petermalina-dev")
	viper.SetDefault("log.level", zap.DebugLevel)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	config := zapdriver.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.Level(viper.GetInt("log.level")))
	logger, err := config.Build(zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName("petermalina"),
	))
	if err != nil {
		panic(err)
	}

	// Listener for gRPC
	lis, err := net.Listen("tcp", ":"+viper.GetString("port"))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	interceptors := grpc_middleware.ChainUnaryServer(
		grpc_zap.UnaryServerInterceptor(logger),
	)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors),
	)

	fbApp, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: viper.GetString("project.id"),
	})

	if err != nil {
		logger.Fatal("firebase app", zap.Error(err))
	}

	fs, err := fbApp.Firestore(ctx)
	if err != nil {
		logger.Fatal("firestore create", zap.Error(err))
	}

	userRepo := user.NewRepository(fs.Collection(internal.CollectionUsers))
	userSvc := user.NewService(logger, userRepo)

	v1.RegisterUserServiceServer(grpcServer, userSvc)
	reflection.Register(grpcServer)
	grpcWebServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)

	srv := &http.Server{
		Addr:    ":" + viper.GetString("port"),
		Handler: multiplexGRPCTraffic(grpcServer, grpcWebServer),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Warn("shutting down grpcServer")
		grpcServer.GracefulStop()
		_ = srv.Close()
	}()

	logger.Info("starting grpcServer", zap.String("port", viper.GetString("port")))
	if err := srv.Serve(lis); err != nil {
		logger.Info("grpcServer exit", zap.Error(err))
	}

}

func multiplexGRPCTraffic(grpcServer *grpc.Server, grpcWebServer *grpcweb.WrappedGrpcServer) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if grpcWebServer.IsAcceptableGrpcCorsRequest(r) || grpcWebServer.IsGrpcWebRequest(r) {
			grpcWebServer.ServeHTTP(w, r)
			return
		}
		grpcServer.ServeHTTP(w, r)
	}), &http2.Server{})
}
