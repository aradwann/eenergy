package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/gapi"
	"github.com/aradwann/eenergy/pb"
	"github.com/aradwann/eenergy/util"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed doc/swagger/*
var content embed.FS

func main() {

	config, err := util.LoadConfig(".", "app")
	if err != nil {
		slog.Error("error loading config: %v\n", err)
		os.Exit(1)
	}

	var logHandler slog.Handler
	if config.Environment == "development" {
		// use simple text output for development
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     nil,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove time.
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{}
				}
				return a
			},
		})
	} else {
		// use JSON Log Handler for production
		logHandler = slog.NewJSONHandler(os.Stdout, nil)

	}
	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		slog.Error("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	db.RunDBMigrations(dbConn, config.MigrationsURL)

	store := db.NewStore(dbConn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		slog.Error("cannot create server: %s", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	// gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterEenergyServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		slog.Error("cannot create listener: %s", err)
	}

	slog.Info(fmt.Sprintf("start gRPC server at %s", listener.Addr().String()))
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("cannot start gRPC server: %s", err)
	}
}

// runGateServer uses in-process translation avoiding extra network hopping, but only allow unary grpc request/response
func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		slog.Error("cannot create server: %s", err)
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
		slog.Error("cannot register handler server: %s", err)
	}

	mux := http.NewServeMux()
	// Handle gRPC requests
	mux.Handle("/", grpcMux)
	// TODO: fix serving URL
	mux.Handle("/swagger/", http.FileServer(http.FS(content)))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		slog.Error("cannot create listener: %s", err)
	}

	slog.Info(fmt.Sprintf("start HTTP gateway server at %s", listener.Addr().String()))
	err = http.Serve(listener, mux)
	if err != nil {
		slog.Error("cannot start http gateway server: %s", err)
	}
}
