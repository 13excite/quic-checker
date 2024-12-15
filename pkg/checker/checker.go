package checker

import (
	"context"
	"log"
	"sync"

	"github.com/13excite/quic-checker/pkg/config"
)

type SiteStatus struct {
	URL        string
	StatusCode int
	Err        error
}

// TODO: add cli and prometheus output modes
func ShellSiteStatusChecker(ctx context.Context, wg *sync.WaitGroup, results <-chan *SiteStatus, config *config.Config) {

	for {
		select {
		case result := <-results:
			if result.Err != nil {
				log.Print("HTTP/3 check error on url: ", result.URL, result.Err)

				wg.Done()
				continue
			}
			if result.StatusCode != config.ExpectedStatusCode {
				log.Print("HTTP/3 check error on url: ", result.URL, " expected status code: ", config.ExpectedStatusCode, " got: ", result.StatusCode)
				wg.Done()
				continue
			}
			log.Print("HTTP/3 check success on url: ", result.URL, " status code: ", result.StatusCode)
			wg.Done()
		case <-ctx.Done():
			return
		}
	}
}
