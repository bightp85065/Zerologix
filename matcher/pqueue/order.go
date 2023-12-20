package pqueue

import (
	"fmt"
	"time"

	"github.com/fmstephe/flib/fmath"
	"github.com/fmstephe/flib/fstrconv"
)

type OrderNode struct {
	priceNode node
	guidNode  node
	amount    uint64
	productId int
	PriceType int
	nextFree  *OrderNode
}

func (o *OrderNode) CopyFrom(amount uint64, productId, priceType, price int, create time.Time) {
	o.amount = amount
	o.productId = productId
	o.PriceType = priceType
	timestamp := create.UnixNano()
	o.setup(uint64(price), uint64(fmath.CombineInt32(int32(timestamp), int32(timestamp))))
}

func (o *OrderNode) setup(price, guid uint64) {
	initNode(o, price, &o.priceNode, &o.guidNode)
	initNode(o, guid, &o.guidNode, &o.priceNode)
}

func (o *OrderNode) Price() uint64 {
	return o.priceNode.val
}

func (o *OrderNode) Guid() uint64 {
	return o.guidNode.val
}

func (o *OrderNode) TraderId() uint32 {
	return uint32(fmath.HighInt32(int64(o.guidNode.val)))
}

func (o *OrderNode) TradeId() uint32 {
	return uint32(fmath.LowInt32(int64(o.guidNode.val)))
}

func (o *OrderNode) Amount() uint64 {
	return o.amount
}

func (o *OrderNode) ReduceAmount(s uint64) {
	o.amount -= s
}

func (o *OrderNode) ProductId() int {
	return o.productId
}

func (o *OrderNode) Remove() {
	o.priceNode.pop()
	o.guidNode.pop()
}

func (o *OrderNode) String() string {
	if o == nil {
		return "<nil>"
	}
	price := fstrconv.ItoaDelim(int64(o.Price()), ',')
	amount := fstrconv.ItoaDelim(int64(o.Amount()), ',')
	traderId := fstrconv.ItoaDelim(int64(o.TraderId()), '-')
	tradeId := fstrconv.ItoaDelim(int64(o.TradeId()), '-')
	productId := fstrconv.ItoaDelim(int64(o.ProductId()), '-')
	// kind := o.kind
	return fmt.Sprintf("price %s, amount %s, trader %s, trade %s, product %s", price, amount, traderId, tradeId, productId)
}
