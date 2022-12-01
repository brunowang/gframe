package gfhttp

import (
	"github.com/brunowang/gframe/gflog"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"net/http"
)

var NullData = struct{}{}

var pbjs = jsonpb.Marshaler{OrigName: true, EmitDefaults: true, EnumsAsInts: false}

var jspb = jsonpb.Unmarshaler{AllowUnknownFields: true}

func BindJson(ctx *gin.Context, jsonReq interface{}) bool {
	var err error
	if msg, ok := jsonReq.(proto.Message); ok {
		err = jspb.Unmarshal(ctx.Request.Body, msg)
	} else {
		err = ctx.ShouldBindJSON(jsonReq)
	}
	if err != nil {
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{
			"status": status,
			"errmsg": http.StatusText(status),
			"data":   NullData,
		})
		gflog.Error(ctx, "bind json error", zap.Error(err))
		return false
	}
	return true
}
