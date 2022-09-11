package gfrds

import (
	"context"
	"fmt"
	"github.com/brunowang/gframe/gflog"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
)

type RedisStream struct {
	cli *redis.Client
}

func NewRedisStream(cli *redis.Client) *RedisStream {
	return &RedisStream{cli: cli}
}

type consumeFunc func(ctx context.Context, msg redis.XMessage) error

func (r *RedisStream) Consume(ctx context.Context, stream, group string, fn consumeFunc, opts ...ConsumeOption) error {
	var lastErr error
	defer func() {
		if err := recover(); err != nil {
			gflog.Error(ctx, "mq consume recovered", zap.String("stream", stream), zap.Any("error", err))
			lastErr = fmt.Errorf("%+v", err)
			return
		}
	}()

	options := consumeOptions{
		ReadCount: 1,
	}
	for _, opt := range opts {
		opt(&options)
	}

	xargs := &redis.XReadGroupArgs{
		Group:    group,
		Consumer: options.ConsumerName,
		Streams:  []string{stream, ">"},
		Count:    options.ReadCount,
	}

	for {
		res, err := r.cli.XReadGroup(xargs).Result()
		if err == redis.ErrClosed {
			gflog.Error(ctx, "mq conn closed", zap.String("stream", stream), zap.Error(err))
			break
		} else if err != nil {
			gflog.Error(ctx, "mq consume failed", zap.String("stream", stream), zap.Error(err))
			continue
		}
		for _, s := range res {
			for _, msg := range s.Messages {
				if err := fn(ctx, msg); err != nil {
					gflog.Error(ctx, "fn invoke failed", zap.String("stream", stream), zap.Error(err))
					continue
				}
				r.cli.XAck(stream, group, msg.ID)
			}
		}
	}
	return lastErr
}

func (r *RedisStream) Produce(ctx context.Context, stream string, msg map[string]interface{}) error {
	if err := r.cli.XAdd(&redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: msg,
	}).Err(); err != nil {
		gflog.Error(ctx, "mq produce failed", zap.String("stream", stream), zap.Error(err))
		return err
	}
	return nil
}
