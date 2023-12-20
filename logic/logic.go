package logic

import (
	"context"
	"zerologix/dao/mysql"
)

//go:generate mockery --name=Logic --inpackage --case underscore
type Logic interface {
	HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error)
}

type Handler struct {
	wdb mysql.DaoI
	rdb mysql.DaoI
}

// TODO kakfka & redis
func NewLogicHandler(rdb, wdb mysql.DaoI) *Handler {
	return &Handler{
		wdb: rdb,
		rdb: wdb,
	}
}

var _ Logic = (*Handler)(nil)
