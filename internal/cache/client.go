package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
)

type CacheClient interface {
	Get(ctx context.Context, key string, category Category) (string, error)
	Set(ctx context.Context, key string, category Category, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string, category Category) error
}

type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

type Cache struct {
	client redis.UniversalClient
	cfg    *Config
}

func NewCacheClient(cfg *Config) (CacheClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        []string{cfg.Addr},
		Password:     cfg.Password,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Cache{
		client: client,
		cfg:    cfg,
	}, nil
}
