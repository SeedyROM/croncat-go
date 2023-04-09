package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/CronCats/croncat-go/pkg/block_stream"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger *logrus.Entry
}

func (a *App) Register() {
	a.Logger.Info("Registering croncat tasks")
}

func (a *App) Unregister() {
	a.Logger.Info("Unregistering croncat tasks")
}

func (a *App) Run() {
	a.Logger.Info("Running croncat tasks")

	provider0 := block_stream.NewRpcBlockProvider("uni-rpc", "https://uni-rpc.reece.sh")
	provider1 := block_stream.NewRpcBlockProvider("uni-rpc2", "https://uni-rpc.reece.sh")

	// Provider waitgroup and context
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Ctrl+C handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println()
		a.Logger.Info("Shutting down...")
		cancel()
	}()

	// Block stream
	blocks := make(chan interface{}, 128)

	// Run the providers
	provider0.Run(ctx, wg, blocks, a.Logger)
	provider1.Run(ctx, wg, blocks, a.Logger)

	// Read the block stream
	go func() {
		for {
			select {
			case <-ctx.Done():
				a.Logger.Info("Stopping block stream")
			default:
				block := <-blocks
				a.Logger.Infof("Got block: %#v", block)
			}
		}
	}()

	wg.Wait()
}
