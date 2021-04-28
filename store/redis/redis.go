package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type DB struct{
	*redis.Client
}

// Returns new redis client
func Dial(ctx context.Context) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		MinIdleConns: 200,	//cfg.Redis.MinIdleConns,
		PoolSize:     12000,//cfg.Redis.PoolSize,
		PoolTimeout:  240,	//time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     "", 	// no password set
		DB:           0,	//cfg.Redis.DB,      
	})
	
	if _, err := client.Ping(ctx).Result(); err != nil {
		fmt.Println("NOT PING")
		return nil, err
	}

	return &DB{client}, nil
}