package consumer

import (
	"context"
	"encoding/json"

	"zerologix/logger"
	"zerologix/matcher/pqueue"
	kafkaModel "zerologix/model/kafka"
	mysqlModel "zerologix/model/mysql"

	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"

	mcr "zerologix/matcher"
)

type orderConsumer struct {
	M *mcr.M
}

func NewOrderConsume() *orderConsumer {
	slabSize := 1024 * 1024 // TODO check slab
	mcr.NewMatcher(slabSize)
	return &orderConsumer{}
}

func (consumer *orderConsumer) Start(ctx context.Context) {
	logger.Log(ctx).Infof("consumer started")
	// dao kafka TODO
	/*cg :=
	cg.RegisterHandler(consumer.handler)
	cg.Run(ctx)*/
}

func (consumer *orderConsumer) handler(sess sarama.ConsumerGroupSession, msgs <-chan *sarama.ConsumerMessage) error {
	g := new(errgroup.Group)
	g.SetLimit(8)

	for msg := range msgs {
		msg := msg
		g.Go(func() error {
			consumer.orderMessageProcess(sess.Context(), msg)
			sess.MarkMessage(msg, "")
			return nil
		})
	}
	return nil
}

func (consumer *orderConsumer) orderMessageProcess(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		err error
	)

	var orderRawMessage kafkaModel.OrderRawMessage
	if err = json.Unmarshal(msg.Value, &orderRawMessage); err != nil {
		logger.Log(ctx).Errorf("unmarshal frontendsrv format fail. err:%v, msg:%+v", err, msg)
		// TODO metric trace
		return
	}
	slab := pqueue.NewSlab(100)
	on := slab.Malloc()
	on.CopyFrom(uint64(orderRawMessage.Quantity), orderRawMessage.Product, orderRawMessage.PriceType,
		orderRawMessage.Price, orderRawMessage.Created)
	switch orderRawMessage.Action {
	case mysqlModel.OrderActionBuy:
		consumer.M.AddBuy(on)
	case mysqlModel.OrderActionSell:
		consumer.M.AddSell(on)
	}

}
