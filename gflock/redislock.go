package gflock

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/atomic"
	"strconv"
	"sync"
	"time"
)

const (
	lockScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

var (
	ErrLockOccupied = errors.New("lock occupied")
)

type RedisLock struct {
	rds  *redis.Client
	once sync.Once
	dur  atomic.Duration
	key  string
	id   string
	err  error
}

func NewRedisLock(key string, rds *redis.Client, uid string) *RedisLock {
	return &RedisLock{
		rds: rds,
		key: key,
		id:  uid,
	}
}

func (r *RedisLock) SetRdsIfAbsent(rds *redis.Client) Locker {
	if r.rds != nil {
		return r
	}
	r.once.Do(func() {
		r.rds = rds
	})
	return r
}

func (r *RedisLock) Lock(ctx context.Context) Locker {
	if errors.Is(r.err, ErrLockOccupied) {
		r.err = nil
	}
	if r.err != nil {
		return r
	}
	millis := r.dur.Load().Milliseconds()
	const tolerance = 500
	expire := strconv.FormatUint(uint64(millis)+tolerance, 10)
	res, err := r.rds.Eval(ctx, lockScript, []string{r.key}, []string{r.id, expire}).Result()
	if errors.Is(err, redis.Nil) {
		r.err = ErrLockOccupied
		return r
	} else if err != nil {
		r.err = err
		return r
	} else if res == nil {
		r.err = ErrLockOccupied
		return r
	}
	reply, ok := res.(string)
	if !ok || reply != "OK" {
		r.err = ErrLockOccupied
		return r
	}
	return r
}

func (r *RedisLock) Unlock(ctx context.Context) Locker {
	if r.err != nil {
		return r
	}
	res, err := r.rds.Eval(ctx, delScript, []string{r.key}, []string{r.id}).Result()
	if err != nil {
		r.err = err
		return r
	}
	reply, ok := res.(int64)
	if !ok || reply != 1 {
		r.err = ErrLockOccupied
		return r
	}
	return r
}

func (r *RedisLock) Expire(dur time.Duration) Locker {
	r.dur.Store(dur)
	return r
}

func (r *RedisLock) Error() error {
	return r.err
}
