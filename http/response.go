package http

import (
	stdhttp "net/http"
	"strings"
	"zerologix/errordefine"

	"github.com/gin-gonic/gin"
)

type baseResponse[Response any] struct {
	ErrorCode int      `json:"code"`
	Success   bool     `json:"success"`
	Message   string   `json:"msg"`
	Response  Response `json:"data"`
}

func (r *baseResponse[Response]) Header() stdhttp.Header {
	if h, ok := (interface{})(r.Response).(Headerer); ok {
		return h.Header()
	}
	return stdhttp.Header{
		"Content-Type": []string{"application/json"},
	}
}

func (r *baseResponse[Response]) Cookie() []*stdhttp.Cookie {
	if cs, ok := (interface{})(r.Response).(Cookier); ok {
		return cs.Cookie()
	}
	return nil
}

func (r *baseResponse[Response]) Code() int {
	if c, ok := (interface{})(r.Response).(StatusCoder); ok {
		return c.Code()
	}
	return stdhttp.StatusOK
}

func newResponse[Response any](rsp Response) *baseResponse[Response] {
	return &baseResponse[Response]{
		ErrorCode: errordefine.OK,
		Success:   true,
		Message:   errordefine.ERR_MAP[errordefine.OK],
		Response:  rsp,
	}
}

func responseHandler[Response any](gc *gin.Context, rsp *Response) error {
	r := newResponse(rsp)
	headers := r.Header()
	code := r.Code()
	cookies := r.Cookie()

	// Set customize headers for particular API response
	for k, v := range headers {
		gc.Header(k, strings.Join(v, ","))
	}

	// Set customize cookies for particular API response
	for _, v := range cookies {
		gc.SetCookie(v.Name, v.Value, v.MaxAge, v.Path, v.Domain, v.Secure, v.HttpOnly)
	}

	// Assume our API only return JSON format
	gc.JSON(code, r)
	return nil
}
