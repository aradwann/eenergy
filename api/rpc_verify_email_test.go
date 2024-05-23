package api

// import (
// 	"context"
// 	"testing"

// 	"github.com/aradwann/eenergy/pb"
// 	mockdb "github.com/aradwann/eenergy/repository/mock"
// 	db "github.com/aradwann/eenergy/repository/store"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// func TestVerifyEmailAPI(t *testing.T) {
// 	emailId := int64(1341)
// 	validSecretCode := "validSecret123fgwfegrtwtqwtqrgwrgewygrwthyhywtra"
// 	invalidEmailId := int64(0)
// 	invalidSecretCode := ""

// 	testCases := []struct {
// 		name          string
// 		req           *pb.VerifyEmailRequest
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, res *pb.VerifyEmailResponse, err error)
// 	}{
// 		{
// 			name: "Success",
// 			req: &pb.VerifyEmailRequest{
// 				EmailId:    emailId,
// 				SecretCode: validSecretCode,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					VerifyEmailTx(gomock.Any(), gomock.Eq(
// 						db.VerifyEmailTxParams{
// 							EmailId:    emailId,
// 							SecretCode: validSecretCode,
// 						})).
// 					Times(1).
// 					Return(db.VerifyEmailTxResult{
// 						User: db.User{IsEmailVerified: true},
// 					}, nil)
// 			},
// 			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
// 				require.NoError(t, err)
// 				require.NotNil(t, res)
// 				require.True(t, res.IsVerified)
// 			},
// 		},
// 		{
// 			name: "InvalidRequest",
// 			req: &pb.VerifyEmailRequest{
// 				EmailId:    invalidEmailId,
// 				SecretCode: invalidSecretCode,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				// No call expected due to validation failure
// 			},
// 			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.InvalidArgument, st.Code())
// 			},
// 		},
// 		{
// 			name: "StoreError",
// 			req: &pb.VerifyEmailRequest{
// 				EmailId:    emailId,
// 				SecretCode: validSecretCode,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					VerifyEmailTx(gomock.Any(), db.VerifyEmailTxParams{
// 						EmailId:    emailId,
// 						SecretCode: validSecretCode,
// 					}).
// 					Times(1).
// 					Return(db.VerifyEmailTxResult{}, status.Errorf(codes.Internal, "store error"))
// 			},
// 			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.Internal, st.Code())
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockStore := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(mockStore)
// 			server := &Server{store: mockStore}

// 			res, err := server.VerifyEmail(context.Background(), tc.req)
// 			tc.checkResponse(t, res, err)
// 		})
// 	}
// }
