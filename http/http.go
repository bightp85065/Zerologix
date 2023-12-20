package http

import (
	"context"
	stdhttp "net/http"

	"gorm.io/gorm"

	"zerologix/logic"
)

type Headerer interface {
	Header() stdhttp.Header
}

type Cookier interface {
	Cookie() []*stdhttp.Cookie
}

type StatusCoder interface {
	Code() int
}

type HttpService interface {
}

type httpService struct {
	logic logic.Logic
	wdb   *gorm.DB
	rdb   *gorm.DB
}

var _ HttpService = (*httpService)(nil)

func NewHttpService(ctx context.Context, logic logic.Logic, rdb, wdb *gorm.DB) *httpService {
	return &httpService{
		logic: logic,
		rdb:   rdb,
		wdb:   wdb,
	}
}
