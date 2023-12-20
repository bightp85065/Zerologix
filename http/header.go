package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey int

const (
	ContextKeyRequestHost contextKey = iota
	ContextKeyREquestProto
	ContextKeyRequestContentType
	ContextKeyRequestMethod
	ContextKeyRequestURI
	ContextKeyRequestPath
	ContextKeyRequestRemoteAddr
)

func populateRequestContext(ctx context.Context, gc *gin.Context) context.Context {
	req := gc.Request

	for k, v := range map[contextKey]string{
		ContextKeyRequestHost:        req.Host,
		ContextKeyREquestProto:       req.Proto,
		ContextKeyRequestContentType: req.Header.Get("Content-Type"),
		ContextKeyRequestMethod:      req.Method,
		ContextKeyRequestURI:         req.RequestURI,
		ContextKeyRequestPath:        req.URL.Path,
		ContextKeyRequestRemoteAddr:  req.RemoteAddr,
	} {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx
}
