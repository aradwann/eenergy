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

func (server *Server) ListUserAccounts(ctx context.Context, req *pb.ListUserAccountsRequest) (*pb.ListUserAccountsResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.AdminRole, util.UserRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateListUserAccountsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	owner := ""
	if authPayload.Role == util.UserRole {
		owner = authPayload.Username
	} else if authPayload.Role == util.AdminRole {
		owner = req.GetUsername()
	}

	arg := db.ListAccountsParams{
		Owner:  owner,
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}

	accs, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "acc not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user accounts: %s", err)
	}

	rsp := &pb.ListUserAccountsResponse{
		Accounts: convertAccounts(accs),
	}
	return rsp, nil
}

func validateListUserAccountsRequest(req *pb.ListUserAccountsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetUsername() != "" {
		if err := val.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, fieldViolation("username", err))
		}
	}
	if err := val.ValidateIntNotNegative(int64(req.GetLimit())); err != nil {
		violations = append(violations, fieldViolation("limit", err))
	}
	if err := val.ValidateIntNotNegative(int64(req.GetOffset())); err != nil {
		violations = append(violations, fieldViolation("offset", err))
	}

	return
}
