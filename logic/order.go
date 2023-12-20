package logic

import (
	"context"
	"time"
	errdef "zerologix/errordefine"
	"zerologix/httperr"
	mysqlModel "zerologix/model/mysql"
)

type OrderCreateRequest struct {
	Product   int `json:"product" binding:"required"`
	Action    int `json:"action" binding:"required,oneof=1 2"` // 1:buy, 2:sell
	Quantity  int `json:"quantity" binding:"required"`
	PriceType int `json:"price_type" binding:"required,oneof=1 2"` // 1:market, 2:limit
	Price     int `json:"price" binding:"required_if=PriceType 2"`
}

type OrderCreateResponse struct {
	OrderId int       `json:"order_id"`
	Created time.Time `json:"created"`
	Status  int       `json:"status"`
}

func (h *Handler) OrderCreate(ctx context.Context, req *OrderCreateRequest) (*OrderCreateResponse, error) {

	order := mysqlModel.Order{
		Creator:   "uid", // TODO
		Product:   req.Product,
		Action:    req.Action,
		Quantity:  req.Quantity,
		PriceType: req.PriceType,
		Price:     req.Price,
		Status:    mysqlModel.OrderStatusPending,
	}

	if err := h.rdb.OrderCreate(ctx, nil, &order); err != nil {
		return nil, httperr.InternalServerErr(errdef.DBError, err.Error())
	}

	// TODO produce kafka

	return &OrderCreateResponse{
		OrderId: order.ID,
		Created: order.Created,
		Status:  mysqlModel.OrderStatusPending,
	}, nil
}
