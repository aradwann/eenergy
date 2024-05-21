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

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.AdminRole, util.UserRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	violations := validateCreateAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	owner := ""
	if authPayload.Role == util.UserRole {
		owner = authPayload.Username
	} else if authPayload.Role == util.AdminRole {
		owner = req.GetUsername()
	}

	arg := db.CreateAccountParams{
		Owner:   owner,
		Balance: 0,
		Unit:    util.KWH,
	}

	acc, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create account: %s", err)
	}

	rsp := &pb.CreateAccountResponse{
		Account: convertAccount(acc),
	}
	return rsp, nil
}

func validateCreateAccountRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetUsername() != "" {
		if err := val.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, fieldViolation("username", err))
		}
	}
	return
}
