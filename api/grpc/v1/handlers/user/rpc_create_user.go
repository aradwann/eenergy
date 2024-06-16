package user

import (
	"context"
	"log/slog"

	"github.com/aradwann/eenergy/api/grpc/v1/handlers/common"
	"github.com/aradwann/eenergy/service/v1/user"
	"github.com/aradwann/eenergy/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	slog.Info("create user RPC", slog.String("username", req.Username))
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, common.InvalidArgumentError(violations)
	}
	user, err := h.service.CreateUser(ctx, user.CreateUserParams{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	res := &CreateUserResponse{
		User: convertUser(user),
	}
	return res, nil
}

func validateCreateUserRequest(req *CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, common.FieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, common.FieldViolation("password", err))
	}
	if err := validator.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, common.FieldViolation("full_name", err))
	}
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, common.FieldViolation("email", err))
	}
	return
}
