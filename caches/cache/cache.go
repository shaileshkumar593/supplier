package cache

import (
	"context"
	"fmt"
	"time"
)

// Cache interface represents the required methods for the factory
type Cache interface {
	Connect(conn string) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	SetTTL(ctx context.Context, key string, value string, ttl time.Duration) error
	SetNX(ctx context.Context, key string, val string, ttl time.Duration) (res bool, err error)
	Delete(ctx context.Context, keys []string) error
	Exist(ctx context.Context, key string) int64
	SetHash(ctx context.Context, key string, values []interface{}) error
	GetHash(ctx context.Context, key string, field string) (string, error)
	InvalidateHash(ctx context.Context, key string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	SetJSON(ctx context.Context, key string, value interface{}) error
	GetJSON(ctx context.Context, key string) (result interface{}, err error)
}

// Factory represents the implementation of the methods
// which will return the Cache interface
type Factory func() (Cache, error)

var factories = make(map[string]Factory)

// Register is a helper method to store all available factory methods
// to the `factories` variable
func Register(name string, factory Factory) error {
	if nil == factory {
		return fmt.Errorf(
			fmt.Sprintf(
				"Cache factory %s does not exist", name))
	}

	_, registered := factories[name]
	if registered {
		return fmt.Errorf(
			fmt.Sprintf(
				"Cache factory %s already registered", name))
	}

	factories[name] = factory

	return nil
}

// New will call the appropriate factory method
// and creates the instance of the Cache interface
func New(name string) (Cache, error) {
	metricsFactory, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("invalid factory name. %s not found in %v", name, factories)
	}

	// Run the Cache factory
	return metricsFactory()
}
