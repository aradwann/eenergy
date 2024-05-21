package api

import (
	"context"
	"errors"

	"github.com/aradwann/eenergy/pb"
	db "github.com/aradwann/eenergy/repository/store"
	"github.com/aradwann/eenergy/util"
	"github.com/aradwann/eenergy/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) TransferEnergy(ctx context.Context, req *pb.TransferEnergyRequest) (*pb.TransferEnergyResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.AdminRole, util.UserRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	violations := validateTransferEnergyRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	fromAccount, err := server.store.GetAccount(ctx, req.FromAccountId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user account: %s", err)
	}
	if authPayload.Username != fromAccount.Owner {
		return nil, status.Error(codes.PermissionDenied, "cannot transfer from this account")
	}

	arg := db.TransferTxParams{
		FromAccountID: req.GetFromAccountId(),
		ToAccountID:   req.GetToAccountId(),
		Amount:        req.GetAmount(),
	}

	transferRes, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to transfer energy units : %s", err)
	}

	rsp := &pb.TransferEnergyResponse{
		Transfer:    convertTransfer(transferRes.Transfer),
		FromAccount: convertAccount(transferRes.FromAccount),
		FromEntry:   convertEntry(transferRes.FromEntry),
	}
	return rsp, nil
}

func validateTransferEnergyRequest(req *pb.TransferEnergyRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateID(req.GetFromAccountId()); err != nil {
		violations = append(violations, fieldViolation("from_account_id", err))
	}

	if err := val.ValidateID(req.GetToAccountId()); err != nil {
		violations = append(violations, fieldViolation("to_account_id", err))
	}

	if err := val.ValidateAmount(req.GetAmount()); err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	return
}
