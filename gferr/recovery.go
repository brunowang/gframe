package gferr

import (
	"context"
	"fmt"
	"github.com/brunowang/gframe/gflog"
	"runtime/debug"
)

func Recovery() {
	if p := recover(); p != nil {
		msg := fmt.Sprintf("%+v\n%s", p, string(debug.Stack()))
		gflog.Error(context.TODO(), msg)
	}
}
