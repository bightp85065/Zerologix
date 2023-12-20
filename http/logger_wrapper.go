package http

import (
	"context"
	"zerologix/constant"
	"zerologix/logger"
	"zerologix/utility"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func loggerWithTraceId(ctx context.Context, gc *gin.Context) context.Context {
	traceId := utility.First(gc.GetHeader(constant.TraceIdHeader), uuid.New().String())

	return logger.WithContext(ctx, logger.
		FromContext(ctx).
		With("trace-id", traceId),
	)
}
