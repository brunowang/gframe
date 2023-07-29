package gfrds

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brunowang/gframe/gflog"
	"github.com/brunowang/gframe/gfserial"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"reflect"
	"time"
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
		ReadCount:    1,
		ReconnectSec: 5,
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
			time.Sleep(time.Duration(options.ReconnectSec) * time.Second)
			continue
		}
		for _, s := range res {
			for _, msg := range s.Messages {
				if err := fn(ctx, msg); err != nil {
					gflog.Error(ctx, "fn invoke failed", zap.String("stream", stream), zap.Error(err))
					r.cli.XAck(stream, group, msg.ID)
					continue
				}
				r.cli.XAck(stream, group, msg.ID)
			}
		}
	}
	return lastErr
}

func (r *RedisStream) Produce(ctx context.Context, stream string, msg interface{}, opts ...ProduceOption) error {
	if msg == nil {
		return fmt.Errorf("mq produce got nil msg")
	}

	options := produceOptions{
		DefaultKey: "message",
	}
	for _, opt := range opts {
		opt(&options)
	}

	vals := make(map[string]interface{})
	if m, ok := msg.(map[string]interface{}); ok {
		vals = m
	} else if m, ok := msg.(*map[string]interface{}); ok {
		vals = *m
	} else if m, ok := msg.(gfserial.Serializable); ok {
		bs, err := m.Serialize()
		if err != nil {
			return err
		}
		if err := json.Unmarshal(bs, &vals); err != nil {
			vals = map[string]interface{}{options.DefaultKey: bs}
		}
	} else if m, ok := msg.(gfserial.Mapper); ok {
		m, err := m.ToMap()
		if err != nil {
			return err
		}
		vals = m
	} else {
		v := reflect.ValueOf(msg)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Invalid {
			return fmt.Errorf("mq produce got invalid msg %+v, typ %T", msg, msg)
		}
		switch v.Kind() {
		case reflect.Struct:
			vals = structs.Map(v.Interface())
		case reflect.Map:
			keys := v.MapKeys()
			for _, k := range keys {
				key := k.String()
				if k.Kind() != reflect.String {
					key = fmt.Sprintf("%v", k.Interface())
				}
				val := v.MapIndex(k)
				vals[key] = val.Interface()
			}
		default:
			vals = map[string]interface{}{options.DefaultKey: v.Interface()}
		}
	}
	if len(vals) == 0 {
		return fmt.Errorf("mq produce got empty msg %+v, typ %T", msg, msg)
	}
	if err := r.cli.XAdd(&redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: vals,
	}).Err(); err != nil {
		gflog.Error(ctx, "mq produce failed", zap.String("stream", stream), zap.Error(err))
		return err
	}
	return nil
}
