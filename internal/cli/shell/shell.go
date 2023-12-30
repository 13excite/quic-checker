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
		Run: func(cmd *cobra.Command, args []string) {
			conf := &config.Config{}
			conf.Defaults()
			conf.ExpectedStatusCode = 400
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			quicConf := quic.Config{}
			quicConf.MaxIdleTimeout = time.Second * 2

			wg := &sync.WaitGroup{}
			inputsChan := make(chan string, 10)
			resultsChan := make(chan *checker.SiteStatus, 10)

			for _, url := range conf.Urls {
				inputsChan <- url
				wg.Add(1)
				// fmt.Println(url)
			}

			go checker.RequesterWorker(inputsChan, resultsChan, quicConf)
			go checker.ShellSiteStatusChecker(ctx, wg, resultsChan, conf)

			wg.Wait()
		},
	}
	return sehllCmd
}
