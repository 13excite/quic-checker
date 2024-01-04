package checker

import (
	"context"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var wrkseq = 1

// Task type represents a worker task descriptor
type Task struct {
	URL string
	WG  *sync.WaitGroup
}

type H3Worker struct {
	id      int
	ctx     context.Context
	queue   chan *Task
	results chan<- *SiteStatus
	client  *QuicClient
}

// NewWorker creates a new worker
func NewWorker(ctx context.Context, quicConf *quic.Config, queue chan *Task, results chan *SiteStatus) {

	roundTripper := &http3.RoundTripper{
		QuicConfig: quicConf,
	}

	w := &H3Worker{
		id:      wrkseq,
		queue:   queue,
		results: results,
		client:  NewClient(roundTripper, 3),
	}

	wrkseq++
	go w.run()
}

// ID is a worker id getter
func (w *H3Worker) ID() int {
	return w.id
}

func (w *H3Worker) run() {
	for {
		select {
		case task := <-w.queue:
			statusCode, err := w.client.Get(task.URL)
			if err != nil {
				w.results <- &SiteStatus{
					URL:        task.URL,
					StatusCode: statusCode,
					Err:        err,
				}
				if task.WG != nil {
					task.WG.Done()
				}
				// next task
				continue
			}
			w.results <- &SiteStatus{
				URL:        task.URL,
				StatusCode: statusCode,
				Err:        nil,
			}
			if task.WG != nil {
				task.WG.Done()
			}
			// next task
			continue

		case <-w.ctx.Done():
			return
		}
	}
}
