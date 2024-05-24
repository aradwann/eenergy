package user

import (
	"context"
	"log/slog"

	"github.com/aradwann/eenergy/api/grpc/v1/handlers/common"
	userService "github.com/aradwann/eenergy/service/v1/user"
	"github.com/aradwann/eenergy/util"
	"github.com/aradwann/eenergy/validator"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	service userService.UserService
	UnimplementedUserServiceServer
}

func NewUserHandler(service userService.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	slog.Info("create user RPC", slog.String("username", req.Username))
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, common.InvalidArgumentError(violations)
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		// AfterCreate: func(user db.User) error {
		// 	taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}
		// 	opts := []asynq.Option{
		// 		asynq.MaxRetry(10),
		// 		asynq.ProcessIn(10 * time.Second), // make room for the DB to commit the transaction before the task is picked up by the worker, otherwise the worker might not find the record
		// 		asynq.Queue(worker.QueueCritical),
		// 	}
		// 	return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		// },
	}

	result, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username already exsits")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &CreateUserResponse{
		User: convertUser(result.User),
	}
	return rsp, nil

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
