package gapi

import (
	"context"
	"testing"
	"time"

	mockdb "github.com/aradwann/eenergy/db/mock"
	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/pb"
	"github.com/aradwann/eenergy/token"
	"github.com/aradwann/eenergy/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListUserAccountsAPI(t *testing.T) {
	user, _ := randomUser(t, util.UserRole)
	admin, _ := randomUser(t, util.AdminRole)

	accounts := []db.Account{
		{
			ID:        1,
			Owner:     user.Username,
			Balance:   1,
			Unit:      util.KWH,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Owner:     user.Username,
			Balance:   2,
			Unit:      util.KWH,
			CreatedAt: time.Now(),
		},
	}
	testCases := []struct {
		name          string
		req           *pb.ListUserAccountsRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.ListUserAccountsResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ListUserAccountsRequest{
				Limit:  5,
				Offset: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  5,
					Offset: 0,
				}

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListUserAccountsResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Len(t, res.Accounts, 2)
			},
		},
		{
			name: "OK by Admin",
			req: &pb.ListUserAccountsRequest{
				Username: &user.Username,
				Limit:    5,
				Offset:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  5,
					Offset: 0,
				}

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, admin.Username, admin.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListUserAccountsResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Len(t, res.Accounts, 2)
			},
		},
		{
			name: "ExpiredToken",
			req:  &pb.ListUserAccountsRequest{},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user.Username, user.Role, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListUserAccountsResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())

			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)

			res, err := server.ListUserAccounts(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
