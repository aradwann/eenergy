package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/token"
	"github.com/aradwann/eenergy/util"
	"github.com/aradwann/eenergy/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}

func newNewContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	ctx := context.Background()
	accessToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotNil(t, payload)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(ctx, md)
}
