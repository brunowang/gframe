package gfrds

type (
	ConsumeOption func(options *consumeOptions)

	consumeOptions struct {
		ConsumerName string
		ReadCount    int64
		ReconnectSec int32
	}

	ProduceOption func(options *produceOptions)

	produceOptions struct {
		DefaultKey string
	}
)

func WithConsumeName(name string) ConsumeOption {
	return func(options *consumeOptions) {
		if name != "" {
			options.ConsumerName = name
		}
	}
}

func WithReadCount(count int64) ConsumeOption {
	return func(options *consumeOptions) {
		if count > 0 {
			options.ReadCount = count
		}
	}
}

func WithReconnectSec(timeSec int32) ConsumeOption {
	return func(options *consumeOptions) {
		if timeSec > 0 {
			options.ReconnectSec = timeSec
		}
	}
}

func WithDefaultKey(key string) ProduceOption {
	return func(options *produceOptions) {
		if key != "" {
			options.DefaultKey = key
		}
	}
}
