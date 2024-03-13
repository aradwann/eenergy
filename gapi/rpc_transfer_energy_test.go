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

func TestTransferEnergyAPI(t *testing.T) {
	user1, _ := randomUser(t)
	account1 := randomAccount(user1)

	user2, _ := randomUser(t)
	account2 := randomAccount(user2)

	amount := int64(1)
	transferId := int64(1)
	fromEntryId := int64(1)
	toEntryId := int64(2)

	testCases := []struct {
		name          string
		req           *pb.TransferEnergyRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.TransferEnergyResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.TransferEnergyRequest{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        1,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				transferRes := db.TransferTxResult{
					Transfer: db.Transfer{
						ID:            transferId,
						FromAccountID: account1.ID,
						ToAccountID:   account2.ID,
						Amount:        amount,
						CreatedAt:     time.Now(),
					},
					FromAccount: db.Account{
						ID:        account1.ID,
						Owner:     account1.Owner,
						Balance:   account1.Balance - amount,
						Unit:      account1.Unit,
						CreatedAt: account1.CreatedAt,
					},
					ToAccount: db.Account{
						ID:        account2.ID,
						Owner:     account2.Owner,
						Balance:   account2.Balance + amount,
						Unit:      account2.Unit,
						CreatedAt: account2.CreatedAt,
					},
					FromEntry: db.Entry{
						ID:        fromEntryId,
						AccountID: account1.ID,
						Amount:    -amount,
						CreatedAt: time.Now(),
					},
					ToEntry: db.Entry{
						ID:        toEntryId,
						AccountID: account2.ID,
						Amount:    amount,
						CreatedAt: time.Now(),
					},
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(transferRes, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user1.Username, user1.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.TransferEnergyResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				fromAccount := res.GetFromAccount()
				require.Equal(t, account1.ID, fromAccount.Id)
				require.Equal(t, account1.Owner, fromAccount.Owner)
				require.Equal(t, account1.Balance-amount, fromAccount.Balance)
				require.Equal(t, account1.Unit, fromAccount.Unit)

				FromEntry := res.GetFromEntry()
				require.Equal(t, FromEntry.Id, fromEntryId)
				require.Equal(t, FromEntry.AccountId, fromAccount.Id)
				require.Equal(t, FromEntry.Amount, -amount)
				require.NotZero(t, FromEntry.CreatedAt)

				transfer := res.GetTransfer()
				require.Equal(t, transfer.Id, transferId)
				require.Equal(t, transfer.Amount, amount)
				require.Equal(t, transfer.FromAccountId, account1.ID)
				require.Equal(t, transfer.ToAccountId, account2.ID)
				require.NotZero(t, transfer.CreatedAt)
			},
		},
		{
			name: "AccountNotFound",
			req: &pb.TransferEnergyRequest{
				FromAccountId: util.RandomInt(4, 53),
				ToAccountId:   util.RandomInt(4, 53),
				Amount:        amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, db.ErrRecordNotFound)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user1.Username, user1.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.TransferEnergyResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())

			},
		},
		{
			name: "ExpiredToken",
			req: &pb.TransferEnergyRequest{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user1.Username, user1.Role, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.TransferEnergyResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())

			},
		},
		{
			name: "PermissionDenied",
			req: &pb.TransferEnergyRequest{
				FromAccountId: account2.ID,
				ToAccountId:   account1.ID,
				Amount:        amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user1.Username, user1.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.TransferEnergyResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.PermissionDenied, st.Code())

			},
		},
		{
			name: "invalidAmount",
			req: &pb.TransferEnergyRequest{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        -amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newNewContextWithBearerToken(t, tokenMaker, user1.Username, user1.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.TransferEnergyResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

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

			res, err := server.TransferEnergy(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func randomAccount(user db.User) db.Account {
	return db.Account{
		ID:        util.RandomInt(0, 900),
		Owner:     user.Username,
		Balance:   util.RandomAmount(),
		Unit:      util.RandomUnit(),
		CreatedAt: time.Now(),
	}
}
