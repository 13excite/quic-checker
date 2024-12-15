package checker

import (
	"context"
	"log"
	"sync"
)

var (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

type SiteStatus struct {
	URL                string
	StatusCode         int
	ExpectedStatusCode int
	Err                error
}

// TODO: add cli and prometheus output modes
func ShellSiteStatusChecker(ctx context.Context, wg *sync.WaitGroup, results <-chan *SiteStatus) {

	for {
		select {
		case result := <-results:
			if result.Err != nil {
				log.Print(colorRed, "HTTP/3 check error on url: ", result.URL, " error: ", result.Err, colorReset)
				wg.Done()
				continue
			}
			if result.StatusCode != result.ExpectedStatusCode {
				log.Print(colorRed, "HTTP/3 check error on url: ", result.URL, " expected status code: ", result.ExpectedStatusCode, " got: ", result.StatusCode, colorReset)
				wg.Done()
				continue
			}
			log.Print(colorGreen, "HTTP/3 check success on url: ", result.URL, " status code: ", result.StatusCode, colorReset)
			wg.Done()
		case <-ctx.Done():
			return
		}
	}
}
