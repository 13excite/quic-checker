package checker

import (
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

// TODO: implement worker pool
func RequesterWorker(
	inputURLs <-chan string,
	results chan<- *SiteStatus,
	quicConf quic.Config,
) {
	roundTripper := &http3.RoundTripper{
		QuicConfig: &quicConf,
	}
	defer roundTripper.Close()

	quicClient := NewClient(roundTripper, 3)

	for url := range inputURLs {
		quicClient.Get(url, results)
	}
}
