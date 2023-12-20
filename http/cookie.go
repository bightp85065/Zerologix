package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

func populateCookieContext(ctx context.Context, gc *gin.Context) context.Context {
	contextKeys := []string{}

	for _, key := range contextKeys {
		// `ErrNoCookie` means the value of cookie cannot be found with key.
		// In our logic, we verified the value whether empty or not, so we could ignore error here.
		sess, _ := gc.Cookie(key)
		ctx = context.WithValue(ctx, key, sess)
	}

	return ctx
}
