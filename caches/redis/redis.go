package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const (
	// CodeRedis ...
	CodeRedis = "redis"

	// EmptyCache ...
	EmptyCache = "redis: nil"
)

type redisContainer struct {
	client *redis.Client
}

var redisConn map[string]*redis.Client

// Connect ...
func (me *redisContainer) Connect(conn string) error {
	me.client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return nil
}

// Get ...
func (me *redisContainer) Get(ctx context.Context, key string) (string, error) {
	var cacheStr, err = me.client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == EmptyCache {
			err = nil
			return "", err
		} else {
			return "", err
		}
	}
	return cacheStr, err
}

// SetNX ...
func (me *redisContainer) SetNX(ctx context.Context, key string, val string, ttl time.Duration) (res bool, err error) {
	res, err = me.client.SetNX(ctx, key, val, ttl).Result()
	if err != nil {
		if err.Error() == EmptyCache {
			err = nil
		}
	}
	return res, err
}

// Set ...
func (me *redisContainer) Set(ctx context.Context, key string, val string) (err error) {
	err = me.client.Set(ctx, key, val, 0).Err()
	if err != nil {
		if err.Error() == EmptyCache {
			err = nil
		}
	}
	return err
}

// SetTTL ...
func (me *redisContainer) SetTTL(ctx context.Context, key string, val string, ttl time.Duration) (err error) {
	err = me.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		if err.Error() == EmptyCache {
			err = nil
		}
	}
	return err

}

// Delete ...
func (me *redisContainer) Delete(ctx context.Context, keys []string) (err error) {
	return me.client.Del(ctx, keys...).Err()
}

// SetHash ...
func (me *redisContainer) SetHash(ctx context.Context, key string, values []interface{}) (err error) {
	return me.client.HSet(ctx, key, values...).Err()
}

// GetHash ...
func (me *redisContainer) GetHash(ctx context.Context, key string, field string) (string, error) {
	return me.client.HGet(ctx, key, field).Result()
}

// InvalidateHash ...
func (me *redisContainer) InvalidateHash(ctx context.Context, key string) error {
	fields, err := me.client.HKeys(ctx, key).Result()
	if err != nil {
		if err.Error() == EmptyCache {
			err = nil
		}
		return err
	}

	_, err = me.client.HDel(ctx, key, fields...).Result()
	return err
}

// Exist ...0 not exist, 1 exist
func (me *redisContainer) Exist(ctx context.Context, key string) int64 {
	return me.client.Exists(ctx, key).Val()
}

// Keys retrieves all keys matching the given pattern.
func (me *redisContainer) Keys(ctx context.Context, pattern string) ([]string, error) {
	return me.client.Keys(ctx, pattern).Result()
}

// SetJSON stores a JSON object under a specific key in Redis
func (me *redisContainer) SetJSON(ctx context.Context, key string, value interface{}) error {
	// Serialize the value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value to JSON: %w", err)
	}

	// Store the JSON string in Redis
	if err := me.client.Set(ctx, key, jsonData, 0).Err(); err != nil {
		return fmt.Errorf("failed to store JSON in Redis: %w", err)
	}

	return nil
}

// GetJSON retrieves and deserializes a JSON object from a specific key in Redis
func (me *redisContainer) GetJSON(ctx context.Context, key string) (result interface{}, err error) {
	// Retrieve the JSON string from Redis
	jsonData, err := me.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON from Redis: %w", err)
	}

	// Deserialize the JSON string into the result object
	if err = json.Unmarshal([]byte(jsonData), result); err != nil {
		return nil, fmt.Errorf("failed to deserialize JSON: %w", err)
	}

	return result, nil
}

// NewRedis ...
func NewRedis() (cache.Cache, error) {
	fmt.Sprintln("redis connection started ")
	ctx := context.Background()

	if _, ok := redisConn["redis"]; ok {
		r := &redisContainer{}
		r.client = redisConn["redis"]
		return r, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Instance().RedisURL, // replace with your Redis URL  localhost:6079  for redis install in pc
		Password: "",                         // no password set
		DB:       0,                          // use default DB
	})

	// Perform a ping-pong check
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	if redisConn == nil {
		redisConn = make(map[string]*redis.Client)
	}
	fmt.Sprintln("successfully connected ")
	redisConn["redis"] = client
	r := &redisContainer{}
	r.client = client

	return r, nil
}
