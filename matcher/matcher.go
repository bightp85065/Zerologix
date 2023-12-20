package matcher

import (
	"zerologix/matcher/pqueue"

	mysqlModel "zerologix/model/mysql"
)

// ref github.com/fmstephe/matching_engine

type productType = int
type priceType = int

type M struct {
	slab           *pqueue.Slab
	newMatchQueues map[productType]map[priceType]*pqueue.MatchQueues
}

func NewMatcher(slabSize int) *M {
	newMatchQueues := make(map[productType]map[priceType]*pqueue.MatchQueues)
	slab := pqueue.NewSlab(slabSize)
	return &M{slab: slab, newMatchQueues: newMatchQueues}
}

func (m *M) getMatchQueues(productId, price int) *pqueue.MatchQueues {
	q := m.newMatchQueues[productId]
	if q == nil {
		m.newMatchQueues[productId] = map[priceType]*pqueue.MatchQueues{}
	}
	q2 := m.newMatchQueues[productId][price]
	if q2 == nil {
		q2 = &pqueue.MatchQueues{}
		m.newMatchQueues[productId][price] = q2
	}
	return q2
}

func (m *M) AddBuy(b *pqueue.OrderNode) {
	if b.PriceType == mysqlModel.OrderPriceTypeMarket {
		panic("It is illegal to send a buy at market price")
	}
	q := m.getMatchQueues(b.ProductId(), int(b.Price()))
	if !m.fillableBuy(b, q) {
		q.PushBuy(b)
	}
}

func (m *M) AddSell(s *pqueue.OrderNode) {
	q := m.getMatchQueues(s.ProductId(), int(s.Price()))
	if !m.fillableSell(s, q) {
		q.PushSell(s)
	}
}

func (m *M) Cancel(o *pqueue.OrderNode) {
	q := m.getMatchQueues(o.ProductId(), int(o.Price()))
	ro := q.Cancel(o)
	if ro != nil {
		m.completeCancelled(ro)
		m.slab.Free(ro)
	} else {
		m.completeNotCancelled(o)
	}
	m.slab.Free(o)
}

func (m *M) fillableBuy(b *pqueue.OrderNode, q *pqueue.MatchQueues) bool {
	for {
		s := q.PeekSell()
		if s == nil {
			return false
		}
		if b.Price() == s.Price() {
			if b.Amount() > s.Amount() {
				amount := s.Amount()
				price := b.Price()
				s.Remove()
				m.slab.Free(s)
				b.ReduceAmount(amount)
				m.completeTrade(b, s, price, amount)
				continue // The sell has been used up
			}
			if s.Amount() > b.Amount() {
				amount := b.Amount()
				price := b.Price()
				s.ReduceAmount(amount)
				m.completeTrade(b, s, price, amount)
				m.slab.Free(b)
				return true // The buy has been used up
			}
			if s.Amount() == b.Amount() {
				amount := b.Amount()
				price := b.Price()
				m.completeTrade(b, s, price, amount)
				s.Remove()
				m.slab.Free(s)
				m.slab.Free(b)
				return true // The buy and sell have been used up
			}
		} else {
			return false
		}
	}
}

func (m *M) fillableSell(s *pqueue.OrderNode, q *pqueue.MatchQueues) bool {
	for {
		b := q.PeekBuy()
		if b == nil {
			return false
		}
		if b.Price() == s.Price() {
			if b.Amount() > s.Amount() {
				amount := s.Amount()
				price := b.Price()
				b.ReduceAmount(amount)
				m.completeTrade(b, s, price, amount)
				s.Remove()
				m.slab.Free(s)
				return true // The sell has been used up
			}
			if s.Amount() > b.Amount() {
				amount := b.Amount()
				price := b.Price()
				s.ReduceAmount(amount)
				m.completeTrade(b, s, price, amount)
				b.Remove()
				m.slab.Free(b) // The buy has been used up
				continue
			}
			if s.Amount() == b.Amount() {
				amount := b.Amount()
				price := b.Price()
				m.completeTrade(b, s, price, amount)
				b.Remove()
				m.slab.Free(b)
				m.slab.Free(s)
				return true // The sell and buy have been used up
			}
		} else {
			return false
		}
	}
}

func (m *M) completeTrade(b, s *pqueue.OrderNode, price, amount uint64) {
	// TODO complete node
	// m.Out.Write(msg.Message{Kind: brk, Price: price, Amount: amount, TraderId: b.TraderId(), TradeId: b.TradeId(), ProductId: b.ProductId()})
	// m.Out.Write(msg.Message{Kind: srk, Price: price, Amount: amount, TraderId: s.TraderId(), TradeId: s.TradeId(), ProductId: s.ProductId()})

	// produce kafka
}

func (m *M) completeCancelled(c *pqueue.OrderNode) {
	// TODO complete node
	// produce kafka
}

func (m *M) completeNotCancelled(nc *pqueue.OrderNode) {
	// TODO complete node
	// produce kafka
}
