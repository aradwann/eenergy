package account

import (
	"context"

	"github.com/aradwann/eenergy/service/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type AccountHandler struct {
	service service.AccountService
	UnimplementedAccountServiceServer
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	account, err := h.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &GetAccountResponse{
		Account: &Account{
			Id:        account.ID,
			Owner:     account.Owner,
			Balance:   account.Balance,
			Unit:      account.Unit,
			CreatedAt: timestamppb.New(account.CreatedAt),
		},
	}, nil
}
