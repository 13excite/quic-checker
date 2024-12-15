package checker

import (
	"context"

	"github.com/quic-go/quic-go"
)

const (
	// dataQueueSize is the size of the queue for the worker pool
	dataQueueSize = 128
)

type WorkerPool struct {
	workersCount int
	queue        chan *Task
	results      chan *SiteStatus
	workers      []*Worker
}

func NewWorkerPool(ctx context.Context, wcount int, quicConf *quic.Config) WorkerPool {
	p := WorkerPool{
		workers:      make([]*Worker, wcount),
		workersCount: wcount,
		queue:        make(chan *Task, dataQueueSize),
		results:      make(chan *SiteStatus, dataQueueSize),
	}
	for i := 0; i < p.workersCount; i++ {
		p.workers[i] = NewWorker(ctx, quicConf, p.queue, p.results)
	}

	return p
}

// AddTask adds a task to the pool queue and increments the WaitGroup counter
func (p *WorkerPool) AddTask(task *Task) {
	if task.WG != nil {
		task.WG.Add(1)
	}
	p.queue <- task
}

// Results returns a channel with the results channel
func (p WorkerPool) Results() <-chan *SiteStatus {
	return p.results
}
