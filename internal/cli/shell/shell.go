package shell

import (
	"context"
	"sync"
	"time"

	"github.com/13excite/quic-checker/pkg/checker"
	"github.com/13excite/quic-checker/pkg/config"
	"github.com/quic-go/quic-go"

	"github.com/spf13/cobra"
)

// NewShellCommand creates a new shell command for running
// the quic checker in shell mode
func NewShellCommand() *cobra.Command {
	sehllCmd := &cobra.Command{
		Use:   "shell",
		Short: "run quic checker in shell mode",
		Run: func(_ *cobra.Command, _ []string) {
			conf := &config.Config{}
			conf.Defaults()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			quicConf := &quic.Config{}
			quicConf.MaxIdleTimeout = time.Second * 2

			wg := &sync.WaitGroup{}
			wp := checker.NewWorkerPool(ctx, conf.GoroutinesCount, quicConf)

			for _, url := range conf.Urls {
				wp.AddTask(&checker.Task{
					URL:                url.URL,
					ExpectedStatusCode: url.ExpectStatusCode,
					WG:                 wg,
				})
			}
			go checker.ShellSiteStatusChecker(ctx, wg, wp.Results())
			wg.Wait()
		},
	}
	return sehllCmd
}
