package gapi

import (
	"context"
	"errors"

	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/pb"
	"github.com/aradwann/eenergy/util"
	"github.com/aradwann/eenergy/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccountEntries(ctx context.Context, req *pb.ListAccountEntriesRequest) (*pb.ListAccountEntriesResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.AdminRole, util.UserRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	account, err := server.store.GetAccount(ctx, req.AccountId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "acc not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get get account: %s", err)
	}

	if authPayload.Role == util.UserRole && authPayload.Username != account.Owner {
		return nil, status.Error(codes.PermissionDenied, "cannot get other user's account entries")
	}

	violations := validateListAccountEntriesRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.ListEntriesParams{
		AccountID: req.GetAccountId(),
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	}

	entries, err := server.store.ListEntries(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "acc not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user accounts: %s", err)
	}

	rsp := &pb.ListAccountEntriesResponse{
		Entries: convertEntries(entries),
	}
	return rsp, nil
}

func validateListAccountEntriesRequest(req *pb.ListAccountEntriesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateID(int64(req.GetAccountId())); err != nil {
		violations = append(violations, fieldViolation("account_id", err))
	}
	if err := val.ValidateIntNotNegative(int64(req.GetLimit())); err != nil {
		violations = append(violations, fieldViolation("limit", err))
	}
	if err := val.ValidateIntNotNegative(int64(req.GetOffset())); err != nil {
		violations = append(violations, fieldViolation("offset", err))
	}

	return
}
