package cache

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

func Init(options *redis.Options) {
	rdb = redis.NewClient(options)
}

func Redis() *redis.Client {
	return rdb
}
