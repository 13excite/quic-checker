package checker

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/quic-go/quic-go/http3"
	"github.com/stretchr/testify/require"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// TODO: Make test dynamic and try to use the RoundTripFunc from the quic-go library
// TODO: Implement client timeout
func TestQuicClient_Get(t *testing.T) {
	// Create an empty QUIC RoundTripper (it will be replaced with a mock client)
	h3RoundTripper := &http3.RoundTripper{}
	// set the max idle connections to 1
	client := NewClient(h3RoundTripper, 1)
	// Replace the quic httpClient with a mock client
	client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
		// TODO: check req url
		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: io.NopCloser(bytes.NewBufferString(`{"staus":"OK"}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	// Create a channel for receiving site status
	siteStatusChan := make(chan *SiteStatus)

	// Perform a GET request using the QuicClient
	go client.Get("http://test.com", siteStatusChan)

	// Receive the site status from the channel
	siteStatus := <-siteStatusChan

	// Check the site status
	require.Empty(t, siteStatus.Err, "Site error is not nil")

	// Check the site response status code
	require.Equal(t, http.StatusOK, siteStatus.StatusCode, "Unexpected status code")

	require.Equal(t, "http://test.com", siteStatus.URL, "Unexpected URL")
}
