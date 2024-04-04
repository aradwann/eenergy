package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/gapi"
	"github.com/aradwann/eenergy/mail"
	"github.com/aradwann/eenergy/observability"
	"github.com/aradwann/eenergy/pb"
	"github.com/aradwann/eenergy/util"
	"github.com/aradwann/eenergy/worker"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

//go:embed doc/swagger/*
var content embed.FS

func main() {

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	config, err := util.LoadConfig(".", ".env")
	if err != nil {
		handleError("error loading config", err)
	}

	initLogger(config)

	dbConn, err := initDatabaseConn(config)
	if err != nil {
		handleError("Unable to connect to database", err)
	}
	defer dbConn.Close()

	runDBMigrations(dbConn, config.MigrationsURL)

	store := db.NewStore(dbConn)
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	go runTaskProcessor(config, redisOpts, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func initLogger(config util.Config) *slog.Logger {
	var logHandler slog.Handler

	if config.Environment == "development" || config.Environment == "test" {
		logHandler = gapi.NewDevelopmentLoggerHandler()
	} else {
		logHandler = gapi.NewProductionLoggerHandler()
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return logger
}

func initDatabaseConn(config util.Config) (*sql.DB, error) {
	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		handleError("Unable to connect to database", err)
	}
	return dbConn, err
}

func runDBMigrations(dbConn *sql.DB, migrationsURL string) {
	db.RunDBMigrations(dbConn, migrationsURL)
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	cert, err := tls.LoadX509KeyPair(config.ServerCrtPath, config.ServerKeyPath)
	if err != nil {
		handleError("cannot load server key pair", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(config.CACrtPath) // If you have a CA certificate, otherwise skip this part
	if err != nil {
		handleError("cannot read ca certificate", err)
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		handleError("failed to append client certs", nil)
	}

	// Create the TLS credentials for the server
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert, // This requires and verifies client certificate
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	})

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		handleError("cannot create gRPC server", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger, grpc.Creds(creds), grpc.StatsHandler(otelgrpc.NewServerHandler()))

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	pb.RegisterEenergyServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		handleError("cannot create listener", err)
	}

	slog.Info(fmt.Sprintf("start gRPC server at %s with TLS", listener.Addr().String()))
	if err := grpcServer.Serve(listener); err != nil {
		handleError("failed to serve", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		handleError("cannot create HTTP gateway server", err)
	}

	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOpts)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterEenergyServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		handleError("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	// TODO: fix path
	mux.Handle("/swagger/", http.FileServer(http.FS(content)))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		handleError("cannot create listener for HTTP gateway server", err)
	}

	slog.Info(fmt.Sprintf("start HTTP gateway server at %s", listener.Addr().String()))
	handler := otelhttp.NewHandler(mux, "your-service-name")
	err = http.Serve(listener, handler)
	if err != nil {
		handleError("cannot start HTTP gateway server", err)
	}
}

func runTaskProcessor(config util.Config, redisOpts asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	taskProcessor := worker.NewRedisTaskProcessor(redisOpts, store, mailer)
	slog.Info("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		handleError("failed to start task processor", err)
	}

}

func handleError(message string, err error) {
	slog.Error("%s: %v", message, err)
	os.Exit(1)
}
