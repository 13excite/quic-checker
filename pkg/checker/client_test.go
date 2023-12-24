package checker

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

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
		Timeout:   3 * time.Second,
		Transport: fn,
	}
}

// TODO: Make test dynamic and try to use the RoundTripFunc from the quic-go library
// TODO: Implement client timeout
func TestQuicClient_Get(t *testing.T) {
	// Create an empty QUIC RoundTripper (it will be replaced with a mock client)
	h3RoundTripper := &http3.RoundTripper{}

	cases := []struct {
		name               string
		response           *http.Response
		siteURL            string
		expectedSiteURL    string
		expectedStatusCode int
		expectedErr        error
	}{
		{
			name: "200 OK responce code",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":"OK"}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			},
			siteURL:            "http://test200.com",
			expectedSiteURL:    "http://test200.com",
			expectedStatusCode: http.StatusOK,
			expectedErr:        nil,
		},
		{
			name: "500 Internal error",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":"Internal Server Error"}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			},
			siteURL:            "http://test500.com",
			expectedSiteURL:    "http://test500.com",
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        nil,
		},
	}
	for _, tc := range cases {
		// set the max idle connections to 1
		client := NewClient(h3RoundTripper, 1)
		// Replace the quic httpClient with a mock client
		client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
			// Check the request URL
			require.Equal(t, tc.expectedSiteURL, req.URL.String(), "Unexpected URL. Case: %s", tc.name)
			return tc.response
		})
		// Create a channel for receiving site status
		siteStatusChan := make(chan *SiteStatus)
		// Perform a GET request using the QuicClient
		go client.Get(tc.siteURL, siteStatusChan)
		// Receive the site status from the channel
		siteStatus := <-siteStatusChan

		// Check the site status
		require.Equal(t, tc.expectedErr, siteStatus.Err, "Unexpected error. Case: %s", tc.name)

		// Check the site response status code
		require.Equal(t, tc.expectedStatusCode, siteStatus.StatusCode, "Unexpected status code. Case: %s", tc.name)
	}
}

func TestQuicClient_AnyError(t *testing.T) {
	// Create a QuicClient with a mock RoundTripper and a short timeout
	h3RoundTripper := &http3.RoundTripper{}
	client := NewClient(h3RoundTripper, 1)
	client.httpClient = http.DefaultClient
	// Create a channel for receiving site status
	siteStatusChan := make(chan *SiteStatus)
	// Perform a GET request to a non-existent server, causing a lookup error
	go client.Get("http://nonexistent", siteStatusChan)

	// Receive the site status from the channel
	siteStatus := <-siteStatusChan

	// Check the site status
	require.Equal(t, "Get \"http://nonexistent\": dial tcp: lookup nonexistent: no such host", siteStatus.Err.Error())

	// Check the site response status code
	require.Equal(t, -1, siteStatus.StatusCode)
}
