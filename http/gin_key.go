package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

func populateGinKeysContext(ctx context.Context, gc *gin.Context) context.Context {
	contextKeys := []string{}

	for _, key := range contextKeys {
		ctx = context.WithValue(ctx, key, gc.GetString(key))
	}

	return ctx
}
