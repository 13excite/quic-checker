package checker

import (
	"log"
)

type SiteStatus struct {
	URL        string
	StatusCode int
	Err        error
}

// TODO: add cli and prometheus output modes
func siteStatusChecker(results <-chan *SiteStatus) {
	for result := range results {
		if result.Err != nil {
			log.Print("HTTP/3 check error on url: ", result.URL, result.Err)
			continue
		}
	}
}
