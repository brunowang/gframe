package gfcache

import (
	"context"
	"errors"
	"github.com/brunowang/gframe/gfserial"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client  *redis.Client
	timeout time.Duration
	nilFlag string
	nilTime time.Duration
	serial  gfserial.Serializer
}

func NewRedisCache(redis *redis.Client) *RedisCache {
	if redis == nil {
		panic(errors.New("got nil redis client"))
	}
	return &RedisCache{
		client:  redis,
		nilFlag: "<nil>",
		nilTime: 2 * time.Second,
		serial:  gfserial.JsonSerializer{},
	}
}

func (c *RedisCache) Timeout(timeout time.Duration) *RedisCache {
	if timeout < 0 {
		return c
	}
	c.timeout = timeout
	return c
}

func (c *RedisCache) NilFlag(nilFlag string) *RedisCache {
	if nilFlag == "" {
		return c
	}
	c.nilFlag = nilFlag
	return c
}

func (c *RedisCache) NilTime(nilTime time.Duration) *RedisCache {
	if nilTime <= 0 {
		return c
	}
	c.nilTime = nilTime
	return c
}

func (c *RedisCache) Serial(serial gfserial.Serializer) *RedisCache {
	if serial == nil {
		return c
	}
	c.serial = serial
	return c
}

func (c *RedisCache) HasCache(key string) bool {
	return c.client.Exists(context.TODO(), key).Val() > 0
}

func (c *RedisCache) GetCache(key string, data interface{}) error {
	val, err := c.client.Get(context.TODO(), key).Bytes()
	if err == redis.Nil {
		return CacheNotFound
	} else if err != nil {
		return err
	} else if string(val) == c.nilFlag {
		return RecordNotFound
	}
	return c.serial.Deserialize(val, data)
}

func (c *RedisCache) SetCache(key string, data interface{}) error {
	val, tmout := []byte(c.nilFlag), c.nilTime
	if data != nil {
		var err error
		val, err = c.serial.Serialize(data)
		if err != nil {
			return err
		}
		tmout = c.timeout
	}
	return c.client.Set(context.TODO(), key, val, tmout).Err()
}

func (c *RedisCache) DelCache(key string) error {
	return c.client.Del(context.TODO(), key).Err()
}
