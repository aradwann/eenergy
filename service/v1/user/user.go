package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aradwann/eenergy/entities"
	DataRepo "github.com/aradwann/eenergy/repository/postgres/user"
	CacheRepo "github.com/aradwann/eenergy/repository/redis"
	"github.com/aradwann/eenergy/util"
)

type UserService interface {
	CreateUser(ctx context.Context, createUserParams CreateUserParams) (*entities.User, error)
	GetUser(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, username string) (*entities.User, error)
}

type userService struct {
	userRepo  DataRepo.UserRepository
	cacheRepo CacheRepo.CacheRepository
}

func NewUserService(userRepo DataRepo.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		// cacheRepo: cacheRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, username string) (*entities.User, error) {
	// Check cache first
	// account, err := s.cacheRepo.GetAccount(ctx, id)
	// if err == nil && account != nil {
	// 	return account, nil
	// }

	// Fetch from database
	user, err := s.userRepo.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	// Cache the result
	// _ = s.cacheRepo.SetAccount(ctx, account)

	return user, nil
}

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}

func (s *userService) CreateUser(ctx context.Context, createUserParams CreateUserParams) (*entities.User, error) {
	hashedPassword, err := util.HashPassword(createUserParams.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %s", err)
	}
	// arg := db.CreateUserTxParams{
	// 	CreateUserParams: db.CreateUserParams{
	// 		Username:       req.GetUsername(),
	// 		HashedPassword: hashedPassword,
	// 		FullName:       req.GetFullName(),
	// 		Email:          req.GetEmail(),
	// 	},
	// AfterCreate: func(user db.User) error {
	// 	taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}
	// 	opts := []asynq.Option{
	// 		asynq.MaxRetry(10),
	// 		asynq.ProcessIn(10 * time.Second), // make room for the DB to commit the transaction before the task is picked up by the worker, otherwise the worker might not find the record
	// 		asynq.Queue(worker.QueueCritical),
	// 	}
	// 	return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	// },
	// }
	return s.userRepo.CreateUser(ctx, DataRepo.CreateUserParams{
		Username:       createUserParams.Username,
		HashedPassword: hashedPassword,
		FullName:       createUserParams.FullName,
		Email:          createUserParams.Email,
	})
}

func (s *userService) UpdateUser(ctx context.Context, username string) (*entities.User, error) {
	return s.userRepo.UpdateUser(ctx, DataRepo.UpdateUserParams{
		HashedPassword:    sql.NullString{},
		PasswordChangedAt: sql.NullTime{},
		FullName:          sql.NullString{},
		Email:             sql.NullString{},
		Username:          username,
		IsEmailVerified:   sql.NullBool{},
	})
}
