package checker

import (
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/quic-go/quic-go/http3"
)

type QuicClient struct {
	httpClient *http.Client
}

func NewClient(h3RoundTripper *http3.RoundTripper, clientTimeout int) *QuicClient {
	return &QuicClient{
		httpClient: &http.Client{
			Transport: h3RoundTripper,
			Timeout:   time.Duration(clientTimeout) * time.Second,
		},
	}
}

func (c *QuicClient) Get(url string, siteStatus chan<- *SiteStatus) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	res, err := c.httpClient.Do(req)
	if err != nil {
		// QoS of UDP is not good enough, especially for small operators, so we need to retry
		// if it's a Client.Timeout , then try to send one more request
		netErr, isNetErr := err.(net.Error)
		if isNetErr && netErr.Timeout() {
			log.Printf("net.Error with a Timeout occured: %v\n", url)
			res, err := c.httpClient.Do(req)
			// close response body if request was successful
			if err == nil {
				defer res.Body.Close()
			}
			siteStatus <- &SiteStatus{url, res.StatusCode, err}
			runtime.Gosched()
			return
		}
	}
	// close response body
	if err == nil {
		defer res.Body.Close()
	}

	siteStatus <- &SiteStatus{url, res.StatusCode, err}
	runtime.Gosched()
}
