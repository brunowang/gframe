package gfhttp

import (
	"encoding/json"
	"github.com/brunowang/gframe/gferr"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"net/http"
)

type Resp struct {
	ctx    *gin.Context
	status int
	code   int
	msg    string
	data   any
	extra  map[string]any
}

func NewResp(ctx *gin.Context) *Resp {
	return &Resp{
		ctx:    ctx,
		status: http.StatusOK,
		data:   NullData,
	}
}

func (r *Resp) Status(status int) *Resp {
	r.status = status
	return r
}

func (r *Resp) Code(code int) *Resp {
	r.code = code
	return r
}

func (r *Resp) Msg(msg string) *Resp {
	r.msg = msg
	return r
}

func (r *Resp) Data(data any) *Resp {
	r.data = data
	return r
}

func (r *Resp) Extra(extra map[string]any) *Resp {
	r.extra = extra
	return r
}

func (r *Resp) Abort() {
	if pb, ok := r.data.(proto.Message); ok {
		js, _ := pbjs.MarshalToString(pb)
		obj := make(map[string]any)
		_ = json.Unmarshal([]byte(js), &obj)
		r.data = obj
	}
	h := gin.H{
		"code": r.code,
		"msg":  r.msg,
		"data": r.data,
	}
	for k, v := range r.extra {
		h[k] = v
	}
	r.ctx.AbortWithStatusJSON(r.status, h)
}

func (r *Resp) OK(data any) {
	r.Data(data).Abort()
}

func (r *Resp) Err(err error) {
	if err == nil {
		r.Abort()
		return
	}
	defer r.ctx.Error(err)
	if e, ok := err.(gferr.IError); ok {
		r.Code(e.Code()).Msg(e.Msg()).Abort()
		return
	}
	r.Msg(err.Error()).Abort()
}
