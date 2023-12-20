package matcher

import (
	"zerologix/matcher/pqueue"

	mysqlModel "zerologix/model/mysql"
)

// ref github.com/fmstephe/matching_engine

type productType = int
type priceType = int

type M struct {
	matchQueues map[productType]*pqueue.MatchQueues
	slab        *pqueue.Slab
}

func NewMatcher(slabSize int) *M {
	matchQueues := make(map[productType]*pqueue.MatchQueues)
	slab := pqueue.NewSlab(slabSize)
	return &M{slab: slab, matchQueues: matchQueues}
}

func (m *M) getMatchQueues(productId int) *pqueue.MatchQueues {
	q := m.matchQueues[productId]
	if q == nil {
		q = &pqueue.MatchQueues{}
		m.matchQueues[productId] = q
	}
	return q
}

func (m *M) AddBuy(b *pqueue.OrderNode) {
	if b.PriceType == mysqlModel.OrderPriceTypeMarket {
		panic("It is illegal to send a buy at market price")
		// TODO check market logic
	}
	q := m.getMatchQueues(b.ProductId())
	if !m.fillableBuy(b, q) {
		q.PushBuy(b)
	}
}

func (m *M) AddSell(s *pqueue.OrderNode) {
	q := m.getMatchQueues(s.ProductId())
	if !m.fillableSell(s, q) {
		q.PushSell(s)
	}
}

func (m *M) Cancel(o *pqueue.OrderNode) {
	q := m.getMatchQueues(o.ProductId())
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
		// s := q.PeekSell()
		s := q.PeekLimitSell(b.Price())
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
		// b := q.PeekBuy()
		b := q.PeekLimitSell(s.Price())
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
