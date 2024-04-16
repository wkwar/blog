package redis

import (
	"backbend/setting"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

/**
redis初始化 以及 关闭
**/

var (
	ctx = context.Background()
	client *redis.Client
)

//初始化redis

func Init(cfg *setting.RedisConf) (err error) {
	client = redis.NewClient(&redis.Options {
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}) 
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return
	}
	return 
}

func Close() {
	_ = client.Close()
}

func Get(key string) (string, error){
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "0", err
	}
	return val, nil
}

func Set(key string, value interface{}, tim time.Duration) error {
	err := client.Set(ctx, key, value, tim).Err()
	if err != nil {
		panic(err)
	}
	return err
}