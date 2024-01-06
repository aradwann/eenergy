package gapi

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGrpcLogger(t *testing.T) {
	// Mock gRPC handler function
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	// Mock gRPC UnaryServerInfo
	mockInfo := &grpc.UnaryServerInfo{
		FullMethod: "/example.Service/Method",
	}

	// Execute GrpcLogger
	_, err := GrpcLogger(context.Background(), nil, mockInfo, mockHandler)
	require.NoError(t, err)

	// Capture logs during the test
	logs := make([]string, 0)
	captureLogs(func() {
		// Execute GrpcLogger
		_, err := GrpcLogger(context.Background(), nil, mockInfo, mockHandler)
		require.NoError(t, err)

	}, func(log string) {
		logs = append(logs, log)
	})

	// Verify the log output
	assert.Contains(t, logs[0], "received grpc req")
	assert.Contains(t, logs[0], "protocol=grpc")
	assert.Contains(t, logs[0], "method=/example.Service/Method")
	assert.Contains(t, logs[0], "status_code=0") // The default value for codes.OK
	assert.Contains(t, logs[0], "status_text=OK")
	assert.Contains(t, logs[0], "duration=")

}

func TestHttpLogger(t *testing.T) {
	// Mock HTTP handler function
	mockHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("OK"))
	})

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Capture logs during the test
	logs := make([]string, 0)
	captureLogs(func() {
		// Execute HttpLogger
		HttpLogger(mockHandler).ServeHTTP(rec, req)
	}, func(log string) {
		logs = append(logs, log)
	})

	// Verify the log output
	assert.Contains(t, logs[0], "received http req")
	assert.Contains(t, logs[0], "status_code=200")
	assert.Contains(t, logs[0], "status_text=OK")
}

func TestHttpLoggerErr(t *testing.T) {
	// Mock HTTP handler function
	mockHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Error"))
	})

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Capture logs during the test
	logs := make([]string, 0)
	captureLogs(func() {
		// Execute HttpLogger
		HttpLogger(mockHandler).ServeHTTP(rec, req)
	}, func(log string) {
		logs = append(logs, log)
	})

	// Verify the log output using substring matching
	assert.Contains(t, logs[0], "received http req")
	assert.Contains(t, logs[0], "status_code=500")
	assert.Contains(t, logs[0], "status_text=\"Internal Server Error\"")
	assert.Contains(t, logs[0], "body=Error")
}

func TestResponseRecorder(t *testing.T) {
	// Mock ResponseWriter
	mockResponseWriter := httptest.NewRecorder()

	// Create a ResponseRecorder
	rec := &ResponseRecorder{
		ResponseWriter: mockResponseWriter,
		StatusCode:     http.StatusOK,
	}

	// Write some data to ResponseRecorder
	rec.Write([]byte("Test Body"))

	// Verify the StatusCode
	require.Equal(t, rec.StatusCode, http.StatusOK)

	// Verify the captured body
	require.Equal(t, string(rec.Body), "Test Body")

}

// captureLogs captures logs produced during the execution of a function.
func captureLogs(fn func(), logCallback func(string)) {
	originalOutput := log.Writer()
	defer func() { log.SetOutput(originalOutput) }()

	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	log.SetOutput(w)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		logCallback(buf.String())
	}()

	fn()

	w.Close()
	wg.Wait()
}
