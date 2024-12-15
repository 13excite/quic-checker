package checker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockQuicClient struct {
	GetFunc func(url string) (int, error)
}

func (m *MockQuicClient) Get(url string) (int, error) {
	return m.GetFunc(url)
}

func TestWorker_Run(t *testing.T) {
	tests := []struct {
		name               string
		task               *Task
		mockGetFunc        func(url string) (int, error)
		expectedStatusCode int
		expectedErr        error
	}{
		{
			name: "successful request",
			task: &Task{
				URL:                "http://example.com",
				ExpectedStatusCode: 200,
				WG:                 &sync.WaitGroup{},
			},
			mockGetFunc: func(_ string) (int, error) {
				return 200, nil
			},
			expectedStatusCode: 200,
			expectedErr:        nil,
		},
		{
			name: "failed request",
			task: &Task{
				URL:                "http://example.com",
				ExpectedStatusCode: 500,
				WG:                 &sync.WaitGroup{},
			},
			mockGetFunc: func(_ string) (int, error) {
				return 500, errors.New("internal server error")
			},
			expectedStatusCode: 500,
			expectedErr:        errors.New("internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			queue := make(chan *Task, 1)
			results := make(chan *SiteStatus, 1)

			worker := &Worker{
				id:      1,
				ctx:     ctx,
				queue:   queue,
				results: results,
				client: &MockQuicClient{
					GetFunc: tt.mockGetFunc,
				},
			}

			go worker.run()

			tt.task.WG.Add(1)
			queue <- tt.task
			// read from the channel and call Done() to signal the task is done
			tt.task.WG.Done()
			// make sure the worker is done and there are no more locks
			tt.task.WG.Wait()

			select {
			case result := <-results:
				require.Equal(t, tt.task.URL, result.URL)
				require.Equal(t, tt.expectedStatusCode, result.StatusCode)
				require.Equal(t, tt.expectedErr, result.Err)
			// timeout after 1 second for being sure the test is not stuck
			case <-time.After(1 * time.Second):
				t.Fatal("test timed out")
			}
		})
	}
}
