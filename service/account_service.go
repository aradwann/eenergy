package service

import (
	"context"

	DataRepo "github.com/aradwann/eenergy/repository/postgres/account"
	CacheRepo "github.com/aradwann/eenergy/repository/redis"
)

type AccountService interface {
	GetAccount(ctx context.Context, id int64) (*DataRepo.Account, error)
}

type accountService struct {
	accountRepo DataRepo.AccountRepository
	cacheRepo   CacheRepo.CacheRepository
}

func NewAccountService(accountRepo DataRepo.AccountRepository, cacheRepo CacheRepo.CacheRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
		cacheRepo:   cacheRepo,
	}
}

func (s *accountService) GetAccount(ctx context.Context, id int64) (*DataRepo.Account, error) {
	// Check cache first
	// account, err := s.cacheRepo.GetAccount(ctx, id)
	// if err == nil && account != nil {
	// 	return account, nil
	// }

	// Fetch from database
	account, err := s.accountRepo.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	// _ = s.cacheRepo.SetAccount(ctx, account)

	return account, nil
}
