// This package is used to create a ticker that will run the checker every interval.
package ticker

import (
	"context"
	"sync"
	"time"

	"github.com/13excite/quic-checker/pkg/checker"

	"github.com/13excite/quic-checker/pkg/config"
	"github.com/quic-go/quic-go"
)

type Job struct {
	config     *config.Config
	quicConfig *quic.Config
	wp         *checker.WorkerPool
}

func NewJob(config *config.Config) *Job {
	quicConfig := &quic.Config{}
	quicConfig.MaxIdleTimeout = time.Second * 2
	return &Job{config: config, quicConfig: quicConfig}
}

func (j *Job) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * time.Duration(j.config.Interval))
	defer ticker.Stop()

	wp := checker.NewWorkerPool(ctx, j.config.GoroutinesCount, j.quicConfig)
	wg := &sync.WaitGroup{}

	for {
		select {
		case <-ticker.C:
			for _, url := range j.config.Urls {
				wp.AddTask(&checker.Task{
					URL:                url.URL,
					ExpectedStatusCode: url.ExpectStatusCode,
					WG:                 wg,
				})
				go checker.ShellSiteStatusChecker(ctx, wg, wp.Results())
				wg.Wait()
			}
		case <-ctx.Done():
			return nil
		}
	}
}
