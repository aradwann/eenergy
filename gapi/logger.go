package gapi

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	startTime := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(startTime)
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	logLevel := slog.LevelInfo
	var errForLog slog.Attr
	if err != nil {
		logLevel = slog.LevelError
		errForLog = slog.String("error", err.Error())
	}
	slog.LogAttrs(context.Background(),
		logLevel,
		"received grpc req",
		slog.String("protocol", "grpc"),
		slog.String("method", info.FullMethod),
		slog.Int("status_code", int(statusCode)),
		slog.String("status_text", statusCode.String()),
		errForLog,
		slog.Duration("duration", duration),
	)

	return res, err
}
