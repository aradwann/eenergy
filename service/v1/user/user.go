package user

import (
	"context"
	"database/sql"

	DataRepo "github.com/aradwann/eenergy/repository/postgres/user"
	CacheRepo "github.com/aradwann/eenergy/repository/redis"
)

type UserService interface {
	CreateUser(ctx context.Context, username string) (*DataRepo.User, error)
	GetUser(ctx context.Context, username string) (*DataRepo.User, error)
	UpdateUser(ctx context.Context, username string) (*DataRepo.User, error)
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

func (s *userService) GetUser(ctx context.Context, username string) (*DataRepo.User, error) {
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

func (s *userService) CreateUser(ctx context.Context, username string) (*DataRepo.User, error) {
	return s.userRepo.CreateUser(ctx, DataRepo.CreateUserParams{
		Username:       username,
		HashedPassword: "",
		FullName:       "",
		Email:          "",
	})
}
func (s *userService) UpdateUser(ctx context.Context, username string) (*DataRepo.User, error) {
	return s.userRepo.UpdateUser(ctx, DataRepo.UpdateUserParams{
		HashedPassword:    sql.NullString{},
		PasswordChangedAt: sql.NullTime{},
		FullName:          sql.NullString{},
		Email:             sql.NullString{},
		Username:          username,
		IsEmailVerified:   sql.NullBool{},
	})
}
