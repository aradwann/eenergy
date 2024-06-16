package redis

import (
	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	// GetAccount(ctx context.Context, id string) (*service.Account, error)
	// SetAccount(ctx context.Context, account *service.Account) error
}

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) CacheRepository {
	return &cacheRepository{client: client}
}

// func (r *cacheRepository) GetAccount(ctx context.Context, id string) (*service.Account, error) {
// 	result, err := r.client.Get(ctx, id).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var account service.Account
// 	if err := json.Unmarshal([]byte(result), &account); err != nil {
// 		return nil, err
// 	}

// 	return &account, nil
// }

// func (r *cacheRepository) SetAccount(ctx context.Context, account *service.Account) error {
// 	data, err := json.Marshal(account)
// 	if err != nil {
// 		return err
// 	}

// 	return r.client.Set(ctx, account.ID, data, 0).Err()
// }
