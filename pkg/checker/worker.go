package checker

import (
	"context"
	"sync"

	"github.com/13excite/quic-checker/pkg/client"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var wrkseq = 1

// Task type represents a Worker task descriptor
type Task struct {
	URL                string
	ExpectedStatusCode int
	WG                 *sync.WaitGroup
}

type QuicClient interface {
	Get(url string) (statusCode int, err error)
}

// Worker type represents a worker object
type Worker struct {
	id      int
	ctx     context.Context
	queue   chan *Task
	results chan<- *SiteStatus
	client  QuicClient
}

// NewWorker creates a new Worker
func NewWorker(ctx context.Context, quicConf *quic.Config, queue chan *Task, results chan *SiteStatus) *Worker {
	w := &Worker{
		ctx:     ctx,
		id:      wrkseq,
		queue:   queue,
		results: results,
		client: client.NewClient(&http3.RoundTripper{
			QuicConfig: quicConf,
		}, 3),
	}
	wrkseq++
	go w.run()

	return w
}

// ID is a Worker id getter
func (w *Worker) ID() int {
	return w.id
}

func (w *Worker) run() {
	for {
		select {
		case task := <-w.queue:
			statusCode, err := w.client.Get(task.URL)
			if err != nil {
				w.results <- &SiteStatus{
					URL:                task.URL,
					StatusCode:         statusCode,
					ExpectedStatusCode: task.ExpectedStatusCode,
					Err:                err,
				}
				// next task
				continue
			}
			w.results <- &SiteStatus{
				URL:                task.URL,
				StatusCode:         statusCode,
				ExpectedStatusCode: task.ExpectedStatusCode,
				Err:                nil,
			}
			// next task
			continue

		case <-w.ctx.Done():
			return
		}
	}
}
