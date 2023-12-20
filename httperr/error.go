package httperr

import (
	"fmt"
	stdhttp "net/http"
)

type HTTPError struct {
	HttpStatusCode int
	ErrorCode      int
	Err            string
}

func (e *HTTPError) Error() string {
	return e.Err
}

func InvalidParamsErr(errCode int, msg string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HttpStatusCode: stdhttp.StatusBadRequest,
		ErrorCode:      errCode,
		Err:            fmt.Sprintf("invalid parameters: "+msg, args...),
	}
}

func UnprocessableRequestErr(errCode int, msg string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HttpStatusCode: stdhttp.StatusUnprocessableEntity,
		ErrorCode:      errCode,
		Err:            fmt.Sprintf("unprocessable request: "+msg, args...),
	}
}

func InternalServerErr(errCode int, msg string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HttpStatusCode: stdhttp.StatusInternalServerError,
		ErrorCode:      errCode,
		Err:            fmt.Sprintf("internal server error: "+msg, args...),
	}
}

func UnauthorizedErr(errCode int, msg string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HttpStatusCode: stdhttp.StatusUnauthorized,
		ErrorCode:      errCode,
		Err:            fmt.Sprintf("unauthorised error: "+msg, args...),
	}
}

func ForbiddenErr(errCode int, msg string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HttpStatusCode: stdhttp.StatusForbidden,
		ErrorCode:      errCode,
		Err:            fmt.Sprintf("forbidden error: "+msg, args...),
	}
}

func ResourceOccupy(errCode int, msg string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HttpStatusCode: stdhttp.StatusBadRequest,
		ErrorCode:      errCode,
		Err:            fmt.Sprintf("resource occupy error: "+msg, args...),
	}
}
