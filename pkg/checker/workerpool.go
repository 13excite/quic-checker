package checker

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

type Job struct {
}

type Result struct {
}

func New(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount),
		Done:         make(chan struct{}),
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go workerToDo(ctx, &wg, wp.jobs, wp.results)
	}
}

func workerToDo(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
}
