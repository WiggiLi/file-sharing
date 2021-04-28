package redis

import (
	"context"
	"encoding/json"
	"time"
	//"log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/WiggiLi/file-sharing-api/model"
)

// UserRedisRepo ...
type UserRepoRedis struct {
	db *DB
}

// NewUserRepo ...
func NewUserRepo(db *DB) *UserRepoRedis {
	return &UserRepoRedis{db: db}
}

// GetUser retrieves user from Redis
func (repo *UserRepoRedis) GetFileNamesByUserID(ctx context.Context, key string) (*[]model.File, error) {
	userBytes, err := repo.db.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "userRedisRepo.GetUser.redisClient.Get")
	}
	user := &[]model.File{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, errors.Wrap(err, "userRedisRepo.GetUser.json.Unmarshal")
	}

	return user, nil
}

// CreateUser creates user in Redis
func (repo *UserRepoRedis) CreateFileNamesByUserID(ctx context.Context, key string, seconds int, user *[]model.File) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "userRedisRepo.CreateUser.json.Marshal")
	}
	if err = repo.db.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "userRedisRepo.CreateUser.Set")
	}
	return nil
}

// DeleteUser deletes user in Redis
func (repo *UserRepoRedis) DeleteFileNamesByUserID(ctx context.Context, key string) error {
	if err := repo.db.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "userRedisRepo.DeleteUser.Del")
	}
	return nil
}
