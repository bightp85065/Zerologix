package logic

import (
	"fmt"
	"time"

	"context"
)

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	SqlResponseTime   string `json:"mysql,omitempty"`
	RedisResponseTime string `json:"redis,omitempty"`
	NowTime           string `json:"now_time,omitempty"`
}

func (h *Handler) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{
		NowTime: fmt.Sprintf("Now: %v", time.Now()),
	}, nil
}
