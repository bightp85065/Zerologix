package http

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"

	errdef "zerologix/errordefine"
	"zerologix/httperr"
	"zerologix/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/exp/slog"
)

type HandlerFunc[Request, Response any] func(ctx context.Context, req *Request) (*Response, error)
type HandlerFuncWithCtx[Request, Response any] func(ctx context.Context, req *Request) (*Response, context.Context, error)
type PopulateContextFunc func(ctx context.Context, gc *gin.Context) context.Context

func resourceBinding[Request any](gc *gin.Context, req *Request) error {
	// If request didn't carry data in body, when unmarshalling, it would return `io.EOF`.
	if err := gc.ShouldBindJSON(req); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	// [ISSUE]: Gin framework didn't support `uri` and `query` in the same struct
	// If we use `uri` and `form` in the same struct, we CANNOT use `binding` with `form` tag at the same field.
	// References:
	//      https://github.com/gin-gonic/gin/pull/3095
	//      https://github.com/gin-gonic/gin/pull/2812
	if err := gc.ShouldBindUri(req); err != nil {
		return err
	}

	if err := gc.ShouldBindWith(req, binding.Form); err != nil {
		return err
	}

	return nil
}

func errorHandler(gc *gin.Context, err error) {
	// fix record: ctx.Errors can't get error
	if err != nil {
		gc.Error(err)
	}
	if e, ok := err.(*httperr.HTTPError); ok {
		gc.JSON(e.HttpStatusCode, map[string]interface{}{
			"code": e.ErrorCode,
			"msg":  e.Error(),
		})
		return
	}
	gc.JSON(stdhttp.StatusBadRequest, map[string]interface{}{"msg": err.Error()})
}

func populateContext(ctx context.Context, gc *gin.Context, fns []PopulateContextFunc) context.Context {
	for _, fn := range fns {
		ctx = fn(ctx, gc)
	}
	return ctx
}

func HandlerWrapper[Request, Response any](ctx context.Context, handler HandlerFunc[Request, Response]) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var req Request

		ctx := populateContext(ctx, gc, []PopulateContextFunc{
			populateRequestContext,
			populateGinKeysContext,
			populateCookieContext,
		})

		ctx = loggerWithTraceId(ctx, gc)
		logger := logger.FromContext(ctx)

		if err := resourceBinding(gc, &req); err != nil {
			errorHandler(gc, httperr.UnprocessableRequestErr(errdef.ParameterError, err.Error()))
			logger.Error("failed to process http request", slog.String("err", err.Error()))
			return
		}

		rsp, err := handler(ctx, &req)
		if err != nil {
			errorHandler(gc, err)
			logger.Error("failed to handler http request", slog.String("err", err.Error()))
			return
		}

		if err := responseHandler(gc, rsp); err != nil {
			errorHandler(gc, err)
			logger.Error("failed to marshal http response", slog.String("err", err.Error()))
			return
		}
	}
}
