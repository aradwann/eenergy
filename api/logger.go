package api

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"log/slog"
)

// GrpcLogger is a unary server interceptor that logs gRPC requests and responses.
func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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
		"received gRPC request",
		slog.String("protocol", "gRPC"),
		slog.String("method", info.FullMethod),
		slog.Int("status_code", int(statusCode)),
		slog.String("status_text", statusCode.String()),
		errForLog,
		slog.Duration("duration", duration),
	)

	return res, err
}

// HttpLogger is a middleware that logs HTTP requests and responses.
func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logLevel := slog.LevelInfo
		var errForLog slog.Attr

		if rec.StatusCode != http.StatusOK {
			logLevel = slog.LevelError
			errForLog = slog.String("body", string(rec.Body))
		}

		slog.LogAttrs(context.Background(),
			logLevel,
			"received HTTP request",
			slog.String("protocol", "HTTP"),
			slog.String("method", req.Method),
			slog.String("uri", req.RequestURI),
			slog.Int("status_code", rec.StatusCode),
			slog.String("status_text", http.StatusText(rec.StatusCode)),
			errForLog,
			slog.Duration("duration", duration),
		)
	})
}

// ResponseRecorder is used to capture the status code and response body.
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// WriteHeader captures the status code.
func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response body.
func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}
