package api

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/aradwann/eenergy/pb"
// 	mockdb "github.com/aradwann/eenergy/repository/mock"
// 	db "github.com/aradwann/eenergy/repository/store"
// 	"github.com/aradwann/eenergy/token"
// 	"github.com/aradwann/eenergy/util"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// func TestListAccountEntriesAPI(t *testing.T) {
// 	user, _ := randomUser(t, util.UserRole)
// 	admin, _ := randomUser(t, util.AdminRole)
// 	account := randomAccount(user)

// 	entries := []db.Entry{
// 		{
// 			ID:        1,
// 			AccountID: account.ID,
// 			Amount:    util.RandomAmount(),
// 			CreatedAt: time.Now(),
// 		},
// 		{
// 			ID:        2,
// 			AccountID: account.ID,
// 			Amount:    util.RandomAmount(),
// 			CreatedAt: time.Now(),
// 		},
// 	}
// 	testCases := []struct {
// 		name          string
// 		req           *pb.ListAccountEntriesRequest
// 		buildStubs    func(store *mockdb.MockStore)
// 		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
// 		checkResponse func(t *testing.T, res *pb.ListAccountEntriesResponse, err error)
// 	}{
// 		{
// 			name: "OK",
// 			req: &pb.ListAccountEntriesRequest{
// 				AccountId: account.ID,
// 				Limit:     5,
// 				Offset:    0,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
// 					Times(1).
// 					Return(account, nil)

// 				arg := db.ListEntriesParams{
// 					AccountID: account.ID,
// 					Limit:     5,
// 					Offset:    0,
// 				}
// 				store.EXPECT().
// 					ListEntries(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(entries, nil)
// 			},
// 			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
// 				return newNewContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, res *pb.ListAccountEntriesResponse, err error) {
// 				require.NoError(t, err)
// 				require.NotNil(t, res)

// 				require.Len(t, res.Entries, 2)
// 			},
// 		},
// 		{
// 			name: "OK by Admin",
// 			req: &pb.ListAccountEntriesRequest{
// 				AccountId: account.ID,
// 				Limit:     5,
// 				Offset:    0,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
// 					Times(1).
// 					Return(account, nil)

// 				arg := db.ListEntriesParams{
// 					AccountID: account.ID,
// 					Limit:     5,
// 					Offset:    0,
// 				}

// 				store.EXPECT().
// 					ListEntries(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(entries, nil)
// 			},
// 			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
// 				return newNewContextWithBearerToken(t, tokenMaker, admin.Username, admin.Role, time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, res *pb.ListAccountEntriesResponse, err error) {
// 				require.NoError(t, err)
// 				require.NotNil(t, res)

// 				require.Len(t, res.Entries, 2)
// 			},
// 		},
// 		{
// 			name: "ExpiredToken",
// 			req:  &pb.ListAccountEntriesRequest{},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					ListEntries(gomock.Any(), gomock.Any()).
// 					Times(0)

// 			},
// 			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
// 				return newNewContextWithBearerToken(t, tokenMaker, user.Username, user.Role, -time.Minute)
// 			},
// 			checkResponse: func(t *testing.T, res *pb.ListAccountEntriesResponse, err error) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.Unauthenticated, st.Code())

// 			},
// 		},
// 	}

// 	for _, tc := range testCases {

// 		t.Run(tc.name, func(t *testing.T) {
// 			storeCtrl := gomock.NewController(t)
// 			store := mockdb.NewMockStore(storeCtrl)

// 			tc.buildStubs(store)
// 			server := newTestServer(t, store, nil)

// 			ctx := tc.buildContext(t, server.tokenMaker)

// 			res, err := server.ListAccountEntries(ctx, tc.req)
// 			tc.checkResponse(t, res, err)
// 		})
// 	}
// }
