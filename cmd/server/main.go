package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	userHandler "github.com/aradwann/eenergy/api/grpc/v1/handlers/user"
	"github.com/aradwann/eenergy/api/grpc/v1/interceptors"
	"github.com/aradwann/eenergy/repository/postgres"
	userRepo "github.com/aradwann/eenergy/repository/postgres/user"
	userService "github.com/aradwann/eenergy/service/v1/user"
	"github.com/aradwann/eenergy/telemetry"

	"github.com/aradwann/eenergy/util"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// main is the entry point of the application.
func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := telemetry.SetupOTelSDK(ctx)
	if err != nil {
		handleError("error setting up OpenTelemetry", err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Load configuration.
	config, err := util.LoadConfig(".", ".env")
	if err != nil {
		handleError("error loading config", err)
	}

	// Initialize logger.
	telemetry.InitLogger(config)

	// Initialize Database
	db := postgres.InitDatabase(config)
	defer db.Close()

	// Set up Redis options for task distribution.
	// redisOpts := asynq.RedisClientOpt{
	// 	Addr: config.RedisAddress,
	// }

	// Run task processor and HTTP gateway server concurrently.
	// taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	// go runTaskProcessor(config, redisOpts, store)
	// go runGatewayServer(config, store, taskDistributor)

	// Initialize repositories
	userRepo := userRepo.NewUserRepository(db)

	// Initialize services
	userService := userService.NewUserService(userRepo)

	// Run gRPC server.
	runGrpcServer(config, userService)
}

// runGrpcServer runs the gRPC server.
func runGrpcServer(config util.Config, userService userService.UserService) {
	// Load TLS certificate and create TLS credentials.
	cert, err := tls.LoadX509KeyPair(config.ServerCrtPath, config.ServerKeyPath)
	if err != nil {
		handleError("cannot load server key pair", err)
	}

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(config.CACrtPath)
	if err != nil {
		handleError("cannot read ca certificate", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		handleError("failed to append client certs", nil)
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	})

	// Create gRPC server.
	// server, err := api.NewServer(config, store, taskDistributor)
	// if err != nil {
	// 	handleError("cannot create gRPC server", err)
	// }

	grpcLogger := grpc.UnaryInterceptor(interceptors.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger, grpc.Creds(creds), grpc.StatsHandler(otelgrpc.NewServerHandler()))

	// Register health check service.
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Register gRPC service.
	userHandler.RegisterUserServiceServer(grpcServer, userHandler.NewUserHandler(userService))

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Start gRPC server.
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		handleError("cannot create listener", err)
	}

	slog.Info(fmt.Sprintf("start gRPC server at %s with TLS", listener.Addr().String()))

	if err := grpcServer.Serve(listener); err != nil {
		handleError("failed to serve", err)
	}
}

// // runGatewayServer runs the HTTP gateway server.
// func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
// 	// Create gRPC server instance.
// 	server, err := api.NewServer(config, store, taskDistributor)
// 	if err != nil {
// 		handleError("cannot create HTTP gateway server", err)
// 	}

// 	// Create JSON marshaller options.
// 	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
// 		MarshalOptions: protojson.MarshalOptions{
// 			UseProtoNames: true,
// 		},
// 		UnmarshalOptions: protojson.UnmarshalOptions{
// 			DiscardUnknown: true,
// 		},
// 	})

// 	// Create gRPC gateway mux.
// 	grpcMux := runtime.NewServeMux(jsonOpts)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Register gRPC server handler.
// 	err = pb.RegisterEenergyServiceHandlerServer(ctx, grpcMux, server)
// 	if err != nil {
// 		handleError("cannot register handler server", err)
// 	}

// 	// Create HTTP mux.
// 	mux := http.NewServeMux()
// 	mux.Handle("/", grpcMux)
// 	// TODO: fix path
// 	mux.Handle("/swagger/", http.FileServer(http.FS(assets.SwaggerFS)))

// 	// Start HTTP gateway server.
// 	listener, err := net.Listen("tcp", config.HTTPServerAddress)
// 	if err != nil {
// 		handleError("cannot create listener for HTTP gateway server", err)
// 	}

// 	slog.Info(fmt.Sprintf("start HTTP gateway server at %s", listener.Addr().String()))

// 	handler := otelhttp.NewHandler(mux, "your-service-name")
// 	err = http.Serve(listener, handler)
// 	if err != nil {
// 		handleError("cannot start HTTP gateway server", err)
// 	}
// }

// // runTaskProcessor runs the task processor.
// func runTaskProcessor(config util.Config, redisOpts asynq.RedisClientOpt, store db.Store) {
// 	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

// 	taskProcessor := worker.NewRedisTaskProcessor(redisOpts, store, mailer)
// 	slog.Info("start task processor")
// 	err := taskProcessor.Start()
// 	if err != nil {
// 		handleError("failed to start task processor", err)
// 	}
// }

// handleError logs an error message and exits the program with status code 1.
func handleError(message string, err error) {
	slog.Error(fmt.Sprintf("%s: %v", message, err))
	os.Exit(1)
}
