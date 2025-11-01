package gfrds

import (
	"context"
	"fmt"
	"github.com/brunowang/gframe/gflog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisPubSub struct {
	cli *redis.Client
}

func NewRedisPubSub(cli *redis.Client) *RedisPubSub {
	return &RedisPubSub{cli: cli}
}

type subscribeFunc func(ctx context.Context, msg redis.Message) error

func (r *RedisPubSub) Subscribe(ctx context.Context, channel string, fn subscribeFunc) error {
	var lastErr error
	defer func() {
		if err := recover(); err != nil {
			gflog.Error(ctx, "mq subscribe recovered", zap.String("channel", channel), zap.Any("error", err))
			lastErr = fmt.Errorf("%+v", err)
			return
		}
	}()

	for msg := range r.cli.Subscribe(channel).Channel() {
		if msg == nil {
			gflog.Error(ctx, "mq subscribe failed", zap.String("channel", channel), zap.Error(fmt.Errorf("got nil message")))
			continue
		}
		if err := fn(ctx, *msg); err != nil {
			gflog.Error(ctx, "fn invoke failed", zap.String("channel", channel), zap.Error(err))
			continue
		}
	}
	return lastErr
}

func (r *RedisPubSub) Publish(ctx context.Context, channel string, msg interface{}) error {
	if err := r.cli.Publish(channel, msg).Err(); err != nil {
		gflog.Error(ctx, "mq publish failed", zap.String("channel", channel), zap.Error(err))
		return err
	}
	return nil
}
