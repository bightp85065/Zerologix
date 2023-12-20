package job

import (
	"context"
	"zerologix/logic"
)

var JobIns *Job

type Job struct{}

func init() {
	JobIns = &Job{}
}

func (j *Job) Run(ctx context.Context) {
	logic.PairJob.Start(ctx)
}

func (j *Job) Close() {
	logic.PairJob.Close()
}
