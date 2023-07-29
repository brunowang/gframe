package gfcache

import (
	"errors"
	"github.com/brunowang/gframe/gfserial"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedigoCache struct {
	client  *redis.Pool
	timeout time.Duration
	nilFlag string
	nilTime time.Duration
	serial  gfserial.Serializer
}

func NewRedigoCache(redis *redis.Pool) *RedigoCache {
	if redis == nil {
		panic(errors.New("got nil redis client"))
	}
	return &RedigoCache{
		client:  redis,
		nilFlag: "<nil>",
		nilTime: 2 * time.Second,
		serial:  gfserial.JsonSerializer{},
	}
}

func (c *RedigoCache) Timeout(timeout time.Duration) *RedigoCache {
	if timeout < 0 {
		return c
	}
	c.timeout = timeout
	return c
}

func (c *RedigoCache) NilFlag(nilFlag string) *RedigoCache {
	if nilFlag == "" {
		return c
	}
	c.nilFlag = nilFlag
	return c
}

func (c *RedigoCache) NilTime(nilTime time.Duration) *RedigoCache {
	if nilTime <= 0 {
		return c
	}
	c.nilTime = nilTime
	return c
}

func (c *RedigoCache) Serial(serial gfserial.Serializer) *RedigoCache {
	if serial == nil {
		return c
	}
	c.serial = serial
	return c
}

func (c *RedigoCache) HasCache(key string) bool {
	conn := c.client.Get()
	defer conn.Close()
	has, _ := redis.Bool(conn.Do("EXISTS", key))
	return has
}

func (c *RedigoCache) GetCache(key string, data interface{}) error {
	conn := c.client.Get()
	defer conn.Close()
	val, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return CacheNotFound
	} else if err != nil {
		return err
	} else if string(val) == c.nilFlag {
		return RecordNotFound
	}
	return c.serial.Deserialize(val, data)
}

func (c *RedigoCache) SetCache(key string, data interface{}) error {
	val, tmout := []byte(c.nilFlag), c.nilTime
	if data != nil {
		var err error
		val, err = c.serial.Serialize(data)
		if err != nil {
			return err
		}
		tmout = c.timeout
	}
	args := []interface{}{key, val}
	if tmout >= time.Millisecond {
		args = append(args, "PX", int64(tmout/time.Millisecond))
	} else if tmout > 0 {
		args = append(args, "PX", 1)
	}
	conn := c.client.Get()
	defer conn.Close()
	_, err := conn.Do("SET", args...)
	return err
}

func (c *RedigoCache) DelCache(key string) error {
	conn := c.client.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
