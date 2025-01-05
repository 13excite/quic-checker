package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/13excite/quic-checker/pkg/config"
	"github.com/13excite/quic-checker/pkg/ticker"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	conf := config.Config{}
	conf.GetConfig(*configPath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	checkTicker := ticker.NewJob(&conf)

	group.Go(func() error {
		return checkTicker.Run(ctx)
	})

	err := group.Wait()
	if err != nil {
		fmt.Println(err)
	}

}
