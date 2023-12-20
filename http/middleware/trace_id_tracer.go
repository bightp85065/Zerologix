package middleware

import (
	"context"
	"fmt"
	"zerologix/constant"
	"zerologix/logger"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func TraceIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(constant.TraceIdHeader)
		if traceId == "" {
			t, _ := uuid.NewV4()
			traceId = t.String()
			ctx.Set(constant.ContextKeyTraceID, traceId)
			logger.Log(ctx).Infof("no x-trace-id found in request %v, generated one : %v", ctx.Request, traceId)
		}
		ctx.Set(constant.ContextKeyTraceID, traceId)
		ctx.Header(constant.TraceIdHeader, traceId)
	}
}

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return "nil-ctx-trace-id"
	}
	if v := ctx.Value(constant.ContextKeyTraceID); v == nil {
		return "local"
	} else {
		return fmt.Sprintf("%v", v)
	}
}
