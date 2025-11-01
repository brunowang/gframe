package gflock

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context) Locker
	Unlock(ctx context.Context) Locker
	Expire(dur time.Duration) Locker
	Error() error
}
