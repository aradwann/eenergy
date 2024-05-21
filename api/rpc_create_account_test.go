package api

import (
	"context"
	"testing"
	"time"

	"github.com/aradwann/eenergy/pb"
	mockdb "github.com/aradwann/eenergy/repository/mock"
	db "github.com/aradwann/eenergy/repository/store"
	"github.com/aradwann/eenergy/token"
	"github.com/aradwann/eenergy/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t, util.UserRole)

	testCases := []struct {
		name          string
		req           *pb.CreateAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateAccountResponse, err error)
	}{
		{
			name: "OK",
			req:  &pb.CreateAccountRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:   user.Username,
					Balance: 0,
					Unit:    util.KWH,
				}
				account := db.Account{
					ID:        1,
					Owner:     user.Username,
					Balance:   0,
					Unit:      util.KWH,
					CreatedAt: time.Now(),
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateAccountResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				require.Equal(t, user.Username, res.Account.Owner)
				require.Equal(t, int64(0), res.Account.Balance)
				require.Equal(t, util.KWH, res.Account.Unit)
				require.NotZero(t, res.Account.CreatedAt)

			},
		},
		{
			name: "ExpiredToken",
			req:  &pb.CreateAccountRequest{},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user.Username, user.Role, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateAccountResponse, err error) {
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

			res, err := server.CreateAccount(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
