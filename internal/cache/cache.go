package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

func (c *Cache) Get(
	ctx context.Context,
	key string,
	category Category,
) (string, error) {
	if err := c.checkInstance(); err != nil {
		return "", err
	}

	key = c.GetKey(key, category)

	res, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%w: error getting cache: %v", ErrInternal, err)
	}

	return res, nil
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	category Category,
	value string,
	ttl time.Duration,
) error {
	if err := c.checkInstance(); err != nil {
		return err
	}

	key = c.GetKey(key, category)

	if err := c.client.Set(
		key,
		value,
		time.Duration(ttl)*time.Second,
	).Err(); err != nil {
		return fmt.Errorf("%w: error setting cache: %v", ErrInternal, err)
	}

	return nil
}

func (c *Cache) Delete(
	ctx context.Context,
	key string,
	category Category,
) error {
	if err := c.checkInstance(); err != nil {
		return err
	}

	key = c.GetKey(key, category)

	if err := c.client.Del(key).Err(); err != nil {
		return fmt.Errorf("%w: error deleting cache: %v", ErrInternal, err)
	}

	return nil
}

func (c *Cache) GetKey(key string, category Category) string {
	return fmt.Sprintf("%s:%s", category, key)
}

func (c *Cache) checkInstance() error {
	if c == nil {
		return ErrInvalidCacheInstance
	}

	if s := c.client.Ping(); s.Err() != nil {
		return fmt.Errorf("%w: %v", ErrInvalidCacheInstance, s.Err().Error())
	}

	return nil
}
