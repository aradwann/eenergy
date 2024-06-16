package user

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aradwann/eenergy/entities"
	DataRepo "github.com/aradwann/eenergy/repository/postgres/user"
	"github.com/aradwann/eenergy/repository/redis"
	"github.com/aradwann/eenergy/util"
)

type UserService interface {
	CreateUser(ctx context.Context, createUserParams CreateUserParams) (*entities.User, error)
	GetUser(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, username string) (*entities.User, error)
}

type userService struct {
	userRepo DataRepo.UserRepository
	jobRepo  redis.JobRepository
	// cacheRepo redis.CacheRepository
	logger *slog.Logger
}

func NewUserService(userRepo DataRepo.UserRepository, jobRepo redis.JobRepository, logger *slog.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		jobRepo:  jobRepo,
		// cacheRepo: cacheRepo,
		logger: logger,
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
	res, err := s.userRepo.CreateUserTx(ctx, DataRepo.CreateUserTxParams{
		CreateUserParams: DataRepo.CreateUserParams{
			Username:       createUserParams.Username,
			HashedPassword: hashedPassword,
			FullName:       createUserParams.FullName,
			Email:          createUserParams.Email,
		},
		AfterCreate: func(user *entities.User) error {
			return s.jobRepo.EnqueueVerificationEmail(ctx, user.Username, user.Email)
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %s", err)
	}

	return res.User, nil
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
