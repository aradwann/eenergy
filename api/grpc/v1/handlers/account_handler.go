package handlers

import (
	"context"

	pb "github.com/aradwann/eenergy/api/grpc/v1/pb"
	"github.com/aradwann/eenergy/service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountHandler struct {
	service service.AccountService
	pb.UnimplementedAccountServiceServer
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := h.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:        account.ID,
			Owner:     account.Owner,
			Balance:   account.Balance,
			Unit:      account.Unit,
			CreatedAt: timestamppb.New(account.CreatedAt),
		},
	}, nil
}
