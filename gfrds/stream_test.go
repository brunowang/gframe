package gfrds

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"testing"
)

func TestRedisStream_Produce(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	cli := redis.NewClient(&redis.Options{Addr: s.Addr()})
	const streamName = "unittest"
	const groupName = "testgroup"
	const defaultKey = "testjob"
	cli.XGroupCreateMkStream(streamName, groupName, "0-0").Err()
	stream := NewRedisStream(cli)

	opt := WithDefaultKey(defaultKey)
	if err := stream.Produce(ctx, streamName, 1234, opt); err != nil {
		panic(err)
	}
	if err := stream.Produce(ctx, streamName, 12.34, opt); err != nil {
		panic(err)
	}
	if err := stream.Produce(ctx, streamName, "abc123", opt); err != nil {
		panic(err)
	}
	stu := struct{ StuName string }{StuName: "teststu"}
	if err := stream.Produce(ctx, streamName, &stu); err != nil {
		panic(err)
	}
	mp := map[string]interface{}{"MapName": "testmap"}
	if err := stream.Produce(ctx, streamName, &mp); err != nil {
		panic(err)
	}

	nilstu := (*struct{ Name string })(nil)
	err = stream.Produce(ctx, streamName, nilstu)
	assert.Error(t, err)
	nilmp := map[string]interface{}(nil)
	err = stream.Produce(ctx, streamName, &nilmp)
	assert.Error(t, err)
	err = stream.Produce(ctx, streamName, nilmp)
	assert.Error(t, err)

	count := atomic.Int32{}
	go stream.Consume(ctx, streamName, groupName, func(ctx context.Context, msg redis.XMessage) error {
		expectKeys := []string{defaultKey, defaultKey, defaultKey, "StuName", "MapName"}
		expectVals := []string{"1234", "12.34", "abc123", "teststu", "testmap"}
		index := int(count.Load())

		for k, v := range msg.Values {
			assert.Equal(t, expectKeys[index], k)
			assert.Equal(t, expectVals[index], v)
		}

		if index == len(expectKeys)-1 {
			cancel()
		}
		count.Add(1)
		return nil
	})

	select {
	case <-ctx.Done():
		fmt.Println("unit test passed!")
	}
}
