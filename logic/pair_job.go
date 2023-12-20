package logic

import (
	"context"
	"time"

	"zerologix/logger"
)

var PairJob = &pairJob{}

type pairJob struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (pr *pairJob) set(ctx context.Context) {
	pr.ctx, pr.cancel = context.WithCancel(ctx)
}

func (pr *pairJob) Start(ctx context.Context) {
	// set init
	pr.set(ctx)

	go func() {
		// todo: use version control
		tr := time.NewTicker(time.Minute * 10)
		for {
			select {
			case <-tr.C:
				pr.do()
			case <-pr.ctx.Done():
				tr.Stop()
				logger.Log(ctx).Infof("insight report job stop.")
				return
			}
		}
	}()
}

func (pr *pairJob) do() {

}

func (pr *pairJob) Close() {
	pr.cancel()
}

func (pr *pairJob) GetName() string {
	return "pair_job"
}
