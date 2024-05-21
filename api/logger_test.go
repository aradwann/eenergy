package api

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Helper function to setup and teardown log capture.
func withCapturedLogs(t *testing.T, fn func()) []string {
	originalOutput := log.Writer()
	defer log.SetOutput(originalOutput)

	r, w, err := os.Pipe()
	require.NoError(t, err)
	log.SetOutput(w)

	var wg sync.WaitGroup
	wg.Add(1)

	var buf bytes.Buffer
	go func() {
		defer wg.Done()
		defer r.Close()
		if _, err := io.Copy(&buf, r); err != nil {
			t.Errorf("Error capturing logs: %v", err)
		}
	}()

	fn()

	require.NoError(t, w.Close())
	wg.Wait()

	return strings.Split(buf.String(), "\n")
}

func TestGrpcLogger(t *testing.T) {
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	mockInfo := &grpc.UnaryServerInfo{FullMethod: "/example.Service/Method"}

	logs := withCapturedLogs(t, func() {
		_, err := GrpcLogger(context.Background(), nil, mockInfo, mockHandler)
		require.NoError(t, err)
	})

	assert.Contains(t, logs[0], "received gRPC request")
	assert.Contains(t, logs[0], "protocol=gRPC")
	assert.Contains(t, logs[0], "method=/example.Service/Method")
	assert.Contains(t, logs[0], "status_code=0") // The default value for codes.OK
	assert.Contains(t, logs[0], "status_text=OK")
	assert.Contains(t, logs[0], "duration=")
}

func TestHttpLogger(t *testing.T) {
	mockHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte("OK"))
		require.NoError(t, err)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	logs := withCapturedLogs(t, func() {
		HttpLogger(mockHandler).ServeHTTP(rec, req)
	})

	assert.Contains(t, logs[0], "received HTTP request")
	assert.Contains(t, logs[0], "status_code=200")
	assert.Contains(t, logs[0], "status_text=OK")
}

func TestHttpLoggerErr(t *testing.T) {
	mockHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusInternalServerError)
		_, err := res.Write([]byte("Error"))
		require.NoError(t, err)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	logs := withCapturedLogs(t, func() {
		HttpLogger(mockHandler).ServeHTTP(rec, req)
	})

	assert.Contains(t, logs[0], "received HTTP request")
	assert.Contains(t, logs[0], "status_code=500")
	assert.Contains(t, logs[0], "status_text=\"Internal Server Error\"")
	assert.Contains(t, logs[0], "body=Error")
}

func TestResponseRecorder(t *testing.T) {
	mockResponseWriter := httptest.NewRecorder()
	rec := &ResponseRecorder{ResponseWriter: mockResponseWriter, StatusCode: http.StatusOK}

	_, err := rec.Write([]byte("Test Body"))
	require.NoError(t, err)

	require.Equal(t, rec.StatusCode, http.StatusOK)
	require.Equal(t, string(rec.Body), "Test Body")
}
