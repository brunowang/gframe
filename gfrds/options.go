package gfrds

type (
	ConsumeOption func(options *consumeOptions)

	consumeOptions struct {
		ConsumerName string
		ReadCount    int64
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
