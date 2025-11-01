package gflock

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var ctx = context.Background()

func runOnRedis(t *testing.T, fn func(client *redis.Client)) {
	srv := miniredis.RunT(t)
	addr := srv.Addr()
	opt := &redis.Options{Addr: addr}
	fn(redis.NewClient(opt))
}

func TestRedisLock(t *testing.T) {
	runOnRedis(t, func(rds *redis.Client) {
		firstLock := NewRedisLock("test", rds, uuid.New().String()).SetRdsIfAbsent(rds).Expire(5 * time.Second)
		assert.Nil(t, firstLock.Error())
		firstLock.Lock(ctx)
		assert.Nil(t, firstLock.Error())

		secondLock := NewRedisLock("test", rds, uuid.New().String()).SetRdsIfAbsent(rds).Expire(5 * time.Second)
		assert.Nil(t, secondLock.Error())
		secondLock.Lock(ctx)
		assert.Equal(t, ErrLockOccupied, secondLock.Error())

		firstLock.Lock(ctx)
		assert.Nil(t, firstLock.Error())
		firstLock.Unlock(ctx)
		assert.Nil(t, firstLock.Error())

		secondLock.Lock(ctx)
		assert.Nil(t, secondLock.Error())
	})
}
