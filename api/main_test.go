package api

import (
	"os"
	"testing"
	"time"

	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// newTestServer create a test server with a config suitable for testing
func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:    util.RandomString(32),
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())

}
