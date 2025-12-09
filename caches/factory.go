package caches

import (
	"swallow-supplier/caches/cache"
	"swallow-supplier/caches/redis"
)

// Init all factories
func Init() {
	cache.Register(redis.CodeRedis, redis.NewRedis)
}
