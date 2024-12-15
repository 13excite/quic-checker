package checker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellSiteStatusChecker(t *testing.T) {

	cases := []struct {
		name               string
		result             *SiteStatus
		expectedLogMessage string
	}{
		{
			name: "Success",
			result: &SiteStatus{
				URL:                "http://example.com",
				StatusCode:         200,
				ExpectedStatusCode: 200,
				Err:                nil,
			},
			expectedLogMessage: "HTTP/3 check success on url: http://example.com status code: 200",
		},
		{
			name: "Error",
			result: &SiteStatus{
				URL:                "http://example.com",
				StatusCode:         0,
				ExpectedStatusCode: 200,
				Err:                fmt.Errorf("network error"),
			},
			expectedLogMessage: "HTTP/3 check error on url: http://example.com error: network error",
		},
		{
			name: "Unexpected status code",
			result: &SiteStatus{
				URL:                "http://example.com",
				StatusCode:         404,
				ExpectedStatusCode: 200,
				Err:                nil,
			},
			expectedLogMessage: "HTTP/3 check error on url: http://example.com expected status code: 200 got: 404",
		},
	}

	for _, tc := range cases {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var wg sync.WaitGroup
		results := make(chan *SiteStatus, 1)

		// Capture log output
		var logOutput strings.Builder
		log.SetOutput(&logOutput)

		go ShellSiteStatusChecker(ctx, &wg, results)

		// Test case: Successful status check
		wg.Add(1)
		results <- tc.result

		// Test case: Status check with an error
		wg.Wait()
		cancel()

		require.Contains(t, logOutput.String(), tc.expectedLogMessage, "Unexpected log message. Case: %s", tc.name)
	}
}
